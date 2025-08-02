# 🚀 Vercel Deployment Guide - Shopping Cart Project

## Overview
This guide covers deploying your fullstack shopping cart application on Vercel with:
- **Frontend**: React + Vite + daisyUI
- **Backend**: Go + Gin + SQLite

## 📋 Prerequisites

1. **GitHub Account**: Create repositories for your code
2. **Vercel Account**: Sign up at [vercel.com](https://vercel.com)
3. **Git**: Installed locally
4. **Two Separate Repositories**: One for frontend, one for backend

## 🏗️ Phase 1: Prepare Your Code

### ✅ Frontend Preparation (Already Done)
- ✅ Environment variable support added to `api.js`
- ✅ `vercel.json` configuration created
- ✅ `.env.example` file created
- ✅ `.gitignore` updated

### ✅ Backend Preparation (Already Done)
- ✅ Environment variable support added to `main.go`
- ✅ Dynamic port configuration
- ✅ Environment-based CORS settings
- ✅ `vercel.json` configuration created
- ✅ `.gitignore` created

## 🔄 Phase 2: Create GitHub Repositories

### 2.1 Create Frontend Repository
```bash
cd "Frontend"
git init
git add .
git commit -m "Initial frontend commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/shopping-cart-frontend.git
git push -u origin main
```

### 2.2 Create Backend Repository
```bash
cd "Backend"
git init
git add .
git commit -m "Initial backend commit"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/shopping-cart-backend.git
git push -u origin main
```

## 🌐 Phase 3: Deploy Backend on Vercel

### 3.1 Deploy Backend
1. Go to [vercel.com](https://vercel.com) and sign in
2. Click "New Project"
3. Import your **backend** repository
4. Vercel will auto-detect it as a Go project
5. **Important**: Set these environment variables in Vercel dashboard:
   ```
   GO_VERSION=1.21
   GIN_MODE=release
   ALLOWED_ORIGINS=*
   ```
6. Click "Deploy"
7. **Note the deployed URL** (e.g., `https://shopping-cart-backend.vercel.app`)

### 3.2 Test Backend Deployment
- Visit: `https://your-backend-url.vercel.app/items`
- Should return JSON response

## 🎨 Phase 4: Deploy Frontend on Vercel

### 4.1 Deploy Frontend
1. In Vercel, click "New Project" again
2. Import your **frontend** repository
3. Vercel will auto-detect it as a Vite project
4. **Important**: Set environment variable:
   ```
   VITE_API_BASE_URL=https://your-backend-url.vercel.app
   ```
   (Replace with your actual backend URL from step 3.1)
5. Click "Deploy"

### 4.2 Configure Build Settings (if needed)
- **Build Command**: `npm run build`
- **Output Directory**: `dist`
- **Install Command**: `npm install`

## 🔧 Phase 5: Post-Deployment Configuration

### 5.1 Update CORS in Backend
After frontend is deployed, update backend environment variables:
```
ALLOWED_ORIGINS=https://your-frontend-url.vercel.app
```

### 5.2 Test Full Application
1. Visit your frontend URL
2. Test login/signup functionality
3. Test adding items to cart
4. Test checkout process
5. Test logout functionality

## 🛠️ Phase 6: Database Considerations

### Current Setup (SQLite)
- ✅ Works for development and demo
- ⚠️ **Note**: Vercel serverless functions are stateless
- ⚠️ Database resets on each deployment

### Production Database Options (Optional)
If you need persistent data:

1. **PlanetScale** (Recommended)
   - MySQL-compatible
   - Serverless
   - Free tier available

2. **Supabase**
   - PostgreSQL
   - Free tier available

3. **MongoDB Atlas**
   - NoSQL option
   - Free tier available

## 📝 Environment Variables Summary

### Frontend (.env in Vercel)
```
VITE_API_BASE_URL=https://your-backend-url.vercel.app
```

### Backend (.env in Vercel)
```
GO_VERSION=1.21
GIN_MODE=release
ALLOWED_ORIGINS=https://your-frontend-url.vercel.app
PORT=8080
```

## 🎯 Quick Deployment Checklist

- [ ] Create two GitHub repositories
- [ ] Push frontend code to frontend repo
- [ ] Push backend code to backend repo
- [ ] Deploy backend on Vercel
- [ ] Note backend URL
- [ ] Deploy frontend on Vercel with backend URL
- [ ] Update CORS settings in backend
- [ ] Test complete application

## 🔗 Useful Commands

### Local Development
```bash
# Frontend
cd Frontend
npm run dev

# Backend
cd Backend
go run main.go
```

### Building for Production
```bash
# Frontend
cd Frontend
npm run build

# Backend
cd Backend
go build -o main .
```

## 📞 Troubleshooting

### Common Issues:

1. **CORS Errors**
   - Ensure `ALLOWED_ORIGINS` is set correctly in backend
   - Check that frontend URL matches exactly

2. **API Connection Failed**
   - Verify `VITE_API_BASE_URL` is set correctly
   - Check backend is deployed and accessible

3. **Database Issues**
   - SQLite works but resets on deployment
   - Consider external database for persistence

4. **Build Failures**
   - Check Node.js version compatibility
   - Verify all dependencies are in package.json

## 🎉 Success!
Your shopping cart application should now be live on Vercel with both frontend and backend deployed separately!

**Frontend URL**: `https://your-frontend-url.vercel.app`
**Backend URL**: `https://your-backend-url.vercel.app`
