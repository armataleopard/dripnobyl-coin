import time
import requests
import os
from datetime import datetime

WEBHOOK_URL = os.getenv('WEBHOOK_URL', 'http://localhost:3000/update-contract')
WEBHOOK_SECRET = os.getenv('WEBHOOK_SECRET', '')
UPDATE_INTERVAL = int(os.getenv('UPDATE_INTERVAL', '300'))  # 5 минут по умолчанию

def update_contract():
    try:
        headers = {
            'Content-Type': 'application/json',
            'X-Webhook-Secret': WEBHOOK_SECRET
        }
        
        response = requests.post(WEBHOOK_URL, headers=headers)
        response.raise_for_status()
        
        print(f'[{datetime.now()}] Contract updated successfully')
        return True
    except Exception as e:
        print(f'[{datetime.now()}] Error updating contract: {str(e)}')
        return False

def main():
    print(f'Starting contract monitor. Update interval: {UPDATE_INTERVAL} seconds')
    
    while True:
        update_contract()
        time.sleep(UPDATE_INTERVAL)

if __name__ == '__main__':
    main() 