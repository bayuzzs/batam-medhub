package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"batam-medhub/internal/auth"
	"batam-medhub/internal/database"
	"batam-medhub/internal/service"
)

func main() {
	fmt.Println("=== 1. Testing JWT Lifetime Validation ===")
	secret := "dev-jwt-signing-secret-minimum-32-bytes-long!"
	issuer := "batam-medhub"
	audience := "batam-medhub-mobile"
	ttl := 15 * time.Minute

	// Valid token
	validToken, err := auth.IssueAccessToken(secret, issuer, audience, "p1", "s1", "SGD", ttl)
	if err != nil {
		panic(err)
	}
	claims, err := auth.ValidateAccessToken(validToken, secret, issuer, audience, ttl)
	if err != nil {
		panic(fmt.Sprintf("valid token failed: %v", err))
	}
	if claims.Subject != "p1" || claims.PreferredCurrency != "SGD" {
		panic("claims mismatch")
	}

	// Expired token (issued in past, exp in past)
	expiredToken, err := auth.IssueAccessToken(secret, issuer, audience, "p1", "s1", "SGD", -1*time.Second)
	if err != nil {
		panic(err)
	}
	_, err = auth.ValidateAccessToken(expiredToken, secret, issuer, audience, ttl)
	if err != auth.ErrTokenExpired && err != auth.ErrInvalidClaims {
		panic(fmt.Sprintf("expected ErrTokenExpired or ErrInvalidClaims, got: %v", err))
	}

	// Mismatched TTL token
	mismatchedToken, err := auth.IssueAccessToken(secret, issuer, audience, "p1", "s1", "SGD", 30*time.Minute)
	if err != nil {
		panic(err)
	}
	_, err = auth.ValidateAccessToken(mismatchedToken, secret, issuer, audience, ttl)
	if err != auth.ErrInvalidClaims {
		panic(fmt.Sprintf("expected ErrInvalidClaims for TTL mismatch, got: %v", err))
	}
	fmt.Println("JWT Lifetime checks PASSED.")

	fmt.Println("\n=== 2. Testing Money & FX Rate Calculations ===")
	dbURL := os.Getenv("DATABASE_URL")
	db, err := database.Open(dbURL)
	if err != nil {
		panic(err)
	}
	moneySvc := service.NewMoneyService(db)
	ctx := context.Background()

	// Negative amount must be rejected
	_, err = moneySvc.Convert(ctx, service.Money{AmountMinor: -500, Currency: "SGD"}, "IDR")
	if err == nil {
		panic("expected error on negative amount, got nil")
	}
	fmt.Printf("Negative amount correctly rejected: %v\n", err)

	// Test conversions and canonicalized rate strings
	// SGD -> IDR
	sgdToIdr, err := moneySvc.Convert(ctx, service.Money{AmountMinor: 10000, Currency: "SGD"}, "IDR")
	if err != nil {
		panic(err)
	}
	fmt.Printf("SGD 100.00 -> IDR %d, FXRate: %s\n", sgdToIdr.Display.AmountMinor, sgdToIdr.FXRate)
	if sgdToIdr.FXRate != "11850.000000" {
		panic(fmt.Sprintf("expected 11850.000000, got %s", sgdToIdr.FXRate))
	}

	// IDR -> SGD
	idrToSgd, err := moneySvc.Convert(ctx, service.Money{AmountMinor: 118500000, Currency: "IDR"}, "SGD")
	if err != nil {
		panic(err)
	}
	fmt.Printf("IDR 1,185,000.00 -> SGD %d, FXRate: %s\n", idrToSgd.Display.AmountMinor, idrToSgd.FXRate)
	if idrToSgd.FXRate != "0.0000843882" {
		panic(fmt.Sprintf("expected 0.0000843882, got %s", idrToSgd.FXRate))
	}

	// Identity SGD -> SGD
	sgdToSgd, err := moneySvc.Convert(ctx, service.Money{AmountMinor: 5000, Currency: "SGD"}, "SGD")
	if err != nil {
		panic(err)
	}
	fmt.Printf("SGD 50.00 -> SGD %d, FXRate: %s\n", sgdToSgd.Display.AmountMinor, sgdToSgd.FXRate)
	if sgdToSgd.FXRate != "1.000000" {
		panic(fmt.Sprintf("expected 1.000000, got %s", sgdToSgd.FXRate))
	}
	fmt.Println("Money & FX Canonicalization checks PASSED.")

	fmt.Println("\n=== ALL DIRECT VERIFICATIONS PASSED ===")
}
