#!/bin/bash

# Deployment Script for Cart Fix

echo "🚀 Deploying Cart Fix to Production..."

echo "📝 Changes made:"
echo "  ✅ Removed hardcoded demo cart items (2 Laptops + 1 Headphones)"
echo "  ✅ New users now start with empty cart"
echo "  ✅ Added debug endpoint POST /carts/clear for testing"

echo ""
echo "🔧 Deploying Backend to Vercel..."
cd Backend
vercel --prod

echo ""
echo "🎨 Deploying Frontend to Vercel..."
cd ../Frontend
npm run build
vercel --prod

echo ""
echo "✅ Deployment completed!"
echo ""
echo "🧪 Testing instructions:"
echo "1. Open your deployed app in incognito/private mode"
echo "2. Login with any credentials" 
echo "3. Cart should be completely empty"
echo "4. Add items - only your added items should appear"
echo "5. Delete items - should work properly"
echo ""
echo "🐛 If issues persist:"
echo "- Clear browser cache completely"
echo "- Try: curl -X POST https://your-backend-url/carts/clear"
echo "- Check browser dev tools console for any errors"

echo ""
echo "🎯 Expected behavior:"
echo "  ✅ Empty cart for new users"
echo "  ✅ No pre-filled demo items"
echo "  ✅ Proper add/delete functionality"
echo "  ✅ User-specific cart isolation"
