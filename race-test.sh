
echo "Initial state:"
curl -s http://localhost:8080/api/vouchers/FLASH50 | jq '{code: .data.code, quota: .data.quota}'
echo ""

echo "Launching 200 concurrent redeems..."
for i in {1..200}; do
    curl -X POST http://localhost:8080/api/vouchers/FLASH50/claim \
         -H "X-User-ID: $i" \
         -s -o /dev/null -w "%{http_code} " &
done | tee /tmp/race_results.txt
wait
echo ""
echo ""

SUCCESS=$(grep -o "200" /tmp/race_results.txt | wc -l | tr -d ' ')
CONFLICT=$(grep -o "409" /tmp/race_results.txt | wc -l | tr -d ' ')

echo "Results:"
echo "   Success (200):  $SUCCESS"
echo "   Conflict (409): $CONFLICT"
echo ""

echo "Final state:"
curl -s http://localhost:8080/api/vouchers/FLASH50 | jq '{code: .data.code, quota: .data.quota}'
echo ""

echo "🔍 Database verification:"
docker exec -it voucher-db psql -U postgres -d voucher-db -c \
    "SELECT code, quota, (SELECT COUNT(*) FROM redemption_history WHERE voucher_id = vouchers.id) as redeemed FROM vouchers WHERE code = 'FLASH50';"

