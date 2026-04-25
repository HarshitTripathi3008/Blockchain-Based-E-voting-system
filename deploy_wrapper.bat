@echo off
ssh -i "C:\Users\harsh\Downloads\chat-app-key.pem" ubuntu@13.62.113.166 "cd MAJOR_PROJECT && git pull origin main && export PATH=$PATH:/usr/local/go/bin && go mod download && go build -o evoting-app main.go && pm2 restart evoting"
