#!/bin/bash

DOMAIN="vc.internetblacksmith.dev"
EXPECTED_IP="161.35.165.206"

echo "=== DNS Propagation Check for $DOMAIN ==="
echo "Expected IP: $EXPECTED_IP"
echo ""

echo "--- Local DNS Resolution ---"
LOCAL_IP=$(dig +short $DOMAIN | tail -n1)
echo "Result: $LOCAL_IP"
if [ "$LOCAL_IP" = "$EXPECTED_IP" ]; then
    echo "✓ Match"
else
    echo "✗ Mismatch"
fi
echo ""

echo "--- Google DNS (8.8.8.8) ---"
GOOGLE_IP=$(dig @8.8.8.8 +short $DOMAIN | tail -n1)
echo "Result: $GOOGLE_IP"
if [ "$GOOGLE_IP" = "$EXPECTED_IP" ]; then
    echo "✓ Match"
else
    echo "✗ Mismatch"
fi
echo ""

echo "--- Cloudflare DNS (1.1.1.1) ---"
CF_IP=$(dig @1.1.1.1 +short $DOMAIN | tail -n1)
echo "Result: $CF_IP"
if [ "$CF_IP" = "$EXPECTED_IP" ]; then
    echo "✓ Match"
else
    echo "✗ Mismatch"
fi
echo ""

echo "--- OpenDNS (208.67.222.222) ---"
OPENDNS_IP=$(dig @208.67.222.222 +short $DOMAIN | tail -n1)
echo "Result: $OPENDNS_IP"
if [ "$OPENDNS_IP" = "$EXPECTED_IP" ]; then
    echo "✓ Match"
else
    echo "✗ Mismatch"
fi
echo ""

echo "--- Summary ---"
if [ "$LOCAL_IP" = "$EXPECTED_IP" ] && [ "$GOOGLE_IP" = "$EXPECTED_IP" ] && [ "$CF_IP" = "$EXPECTED_IP" ] && [ "$OPENDNS_IP" = "$EXPECTED_IP" ]; then
    echo "✓ DNS fully propagated"
else
    echo "✗ DNS propagation incomplete"
fi
