#!/bin/bash

# Проверяем наличие изменений
if [ -z "$(git status --porcelain)" ]; then
    echo "No changes to commit"
    exit 0
fi

# Получаем текущую дату и время
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# Добавляем все изменения
git add .

# Коммитим изменения с временной меткой
git commit -m "Auto-sync update: $TIMESTAMP"

# Пушим изменения в репозиторий
git push origin main || git push origin master

# Выводим статус
echo "Changes synchronized at $TIMESTAMP" 