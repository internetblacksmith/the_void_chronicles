// Copyright (C) 2024 Paolo Fabbri
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"log"
	"net"
	"sync"
	"time"
)

type rateLimitEntry struct {
	attempts     int
	firstAttempt time.Time
	blockedUntil time.Time
}

// RateLimiter tracks connection attempts and enforces rate limiting.
type RateLimiter struct {
	mu      sync.RWMutex
	entries map[string]*rateLimitEntry

	maxAttempts     int
	windowDuration  time.Duration
	blockDuration   time.Duration
	cleanupInterval time.Duration
}

// NewRateLimiter creates a new rate limiter with specified parameters.
func NewRateLimiter(maxAttempts int, windowDuration, blockDuration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries:         make(map[string]*rateLimitEntry),
		maxAttempts:     maxAttempts,
		windowDuration:  windowDuration,
		blockDuration:   blockDuration,
		cleanupInterval: 10 * time.Minute,
	}

	go rl.startCleanup()

	return rl
}

// AllowConnection checks if a connection from the given IP should be allowed.
func (rl *RateLimiter) AllowConnection(addr net.Addr) bool {
	ip := extractIP(addr)
	if ip == "" {
		return true
	}

	rl.mu.RLock()
	entry, exists := rl.entries[ip]
	rl.mu.RUnlock()

	now := time.Now()

	if exists {
		if now.Before(entry.blockedUntil) {
			log.Printf("Rate limit: blocked connection from %s (blocked until %s)", ip, entry.blockedUntil.Format(time.RFC3339))
			return false
		}

		if now.Sub(entry.firstAttempt) > rl.windowDuration {
			rl.mu.Lock()
			entry.attempts = 0
			entry.firstAttempt = now
			rl.mu.Unlock()
		}
	}

	return true
}

// RecordFailedAuth records a failed authentication attempt.
func (rl *RateLimiter) RecordFailedAuth(addr net.Addr) {
	ip := extractIP(addr)
	if ip == "" {
		return
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.entries[ip]

	if !exists {
		entry = &rateLimitEntry{
			attempts:     1,
			firstAttempt: now,
		}
		rl.entries[ip] = entry
	} else {
		if now.Sub(entry.firstAttempt) > rl.windowDuration {
			entry.attempts = 1
			entry.firstAttempt = now
			entry.blockedUntil = time.Time{}
		} else {
			entry.attempts++
		}
	}

	if entry.attempts >= rl.maxAttempts {
		entry.blockedUntil = now.Add(rl.blockDuration)
		log.Printf("Rate limit: IP %s blocked for %v after %d failed attempts", ip, rl.blockDuration, entry.attempts)
	}
}

// RecordSuccessfulAuth resets the failed attempt counter for an IP.
func (rl *RateLimiter) RecordSuccessfulAuth(addr net.Addr) {
	ip := extractIP(addr)
	if ip == "" {
		return
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.entries, ip)
}

func (rl *RateLimiter) startCleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cleaned := 0

	for ip, entry := range rl.entries {
		if now.After(entry.blockedUntil) && now.Sub(entry.firstAttempt) > rl.windowDuration {
			delete(rl.entries, ip)
			cleaned++
		}
	}

	if cleaned > 0 {
		log.Printf("Rate limit: cleaned up %d expired entries", cleaned)
	}
}

func extractIP(addr net.Addr) string {
	if addr == nil {
		return ""
	}

	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return addr.String()
	}

	return host
}
