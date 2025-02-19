#!/bin/bash

# Backup script for minecraft world, Uploads to external storeage and then sends download to Discord Webhook to share with other players
SERVER_DIR="/home/world"
BACKUP_DIR="/home/backup"
UPLOAD_URL="https://bashupload.com/"
WEBHOOK_URL=""

send_webhook() {
  message="$1"
  curl -X POST -H "Content-Type: application/json" -d "{\"content\": \"$message\"}" "$WEBHOOK_URL"
}

timestamp=$(date +%Y-%m-%d_%H-%M)

backup_file="minecraft_world_${timestamp}.tar.gz"

full_backup_path="$BACKUP_DIR/$backup_file"

tar -czvf "$full_backup_path" "$SERVER_DIR"

upload_response=$(curl -s -T "$full_backup_path" "$UPLOAD_URL")

if [[ "$upload_response" == *"https"* ]]; then
  uploaded_url=$(echo "$upload_response" | grep -oP 'https://[^\s]*')

  message="Minecraft world backup uploaded successfully!\n\`\`\`\n$uploaded_url\n\`\`\`"

  send_webhook "$message"
  
  rm "$full_backup_path"
else
  send_webhook "Error uploading Minecraft world backup! Could not retrieve URL."
fi

exit 0
