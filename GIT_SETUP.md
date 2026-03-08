# Git Setup Instructions

## Option 1: Using HTTPS (Recommended if SSH keys not set up)

After creating the repository on GitHub, run:

```powershell
# Add remote (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/lighthouse.git

# Push to GitHub
git branch -M main
git push -u origin main
```

## Option 2: Using SSH (If you have SSH keys set up)

After creating the repository on GitHub, run:

```powershell
# Add remote (replace YOUR_USERNAME with your GitHub username)
git remote add origin git@github.com:YOUR_USERNAME/lighthouse.git

# Push to GitHub
git branch -M main
git push -u origin main
```

## If Repository Already Exists

If the repository `gharibvand/lighthouse` already exists and you have access:

### Using HTTPS:
```powershell
git remote add origin https://github.com/gharibvand/lighthouse.git
git branch -M main
git push -u origin main
```

### Using SSH (if you have SSH keys):
```powershell
git remote add origin git@github.com:gharibvand/lighthouse.git
git branch -M main
git push -u origin main
```

## Setting up SSH Keys (Optional)

If you want to use SSH instead of HTTPS:

1. Generate SSH key:
```powershell
ssh-keygen -t ed25519 -C "your_email@example.com"
```

2. Start ssh-agent:
```powershell
Start-Service ssh-agent
ssh-add ~/.ssh/id_ed25519
```

3. Copy public key:
```powershell
Get-Content ~/.ssh/id_ed25519.pub | Set-Clipboard
```

4. Add to GitHub:
   - Go to https://github.com/settings/keys
   - Click "New SSH key"
   - Paste your key and save
