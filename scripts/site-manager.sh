#!/bin/bash

# Директория для всех сайтов
SITES_DIR="/Users/stanislav.uy/Desktop/sites"
CURRENT_SITE_FILE="$SITES_DIR/current_site.txt"

# Функция для создания нового сайта
create_new_site() {
    # Остановить текущий мониторинг если есть
    pkill -f "contract-monitor"
    
    # Создаем новую директорию для сайта
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    NEW_SITE_DIR="$SITES_DIR/elmofi_$TIMESTAMP"
    
    # Создаем структуру директорий
    mkdir -p "$NEW_SITE_DIR"
    mkdir -p "$NEW_SITE_DIR/images"
    
    # Копируем только файлы для мониторинга контракта
    cp update-contract.sh contract-monitor.go config.json start-monitor.sh "$NEW_SITE_DIR/"
    
    # Копируем новые файлы дизайна
    if [ -f "index.html" ]; then
        cp index.html "$NEW_SITE_DIR/"
    fi
    if [ -f "styles.css" ]; then
        cp styles.css "$NEW_SITE_DIR/"
    fi
    if [ -f "script.js" ]; then
        cp script.js "$NEW_SITE_DIR/"
    fi
    
    # Копируем все изображения, если они есть
    if [ -d "images" ]; then
        cp images/* "$NEW_SITE_DIR/images/" 2>/dev/null || true
    fi
    
    # Копируем все PNG файлы из текущей директории
    cp *.png "$NEW_SITE_DIR/" 2>/dev/null || true
    
    # Обновляем ссылку на текущий сайт
    echo "$NEW_SITE_DIR" > "$CURRENT_SITE_FILE"
    
    echo "=== Создание нового сайта ==="
    echo "Директория: $NEW_SITE_DIR"
    echo "Скопированные файлы:"
    ls -la "$NEW_SITE_DIR"
    echo "Изображения:"
    ls -la "$NEW_SITE_DIR/images" 2>/dev/null || echo "Нет изображений"
    
    cd "$NEW_SITE_DIR"
    
    # Запускаем мониторинг для нового сайта
    ./start-monitor.sh &
    
    echo "=== Мониторинг запущен ==="
    echo "Сайт готов к использованию"
}

# Функция для удаления старого сайта
delete_old_site() {
    if [ -f "$CURRENT_SITE_FILE" ]; then
        OLD_SITE=$(cat "$CURRENT_SITE_FILE")
        if [ -d "$OLD_SITE" ]; then
            echo "=== Удаление старого сайта ==="
            echo "Директория: $OLD_SITE"
            
            # Останавливаем мониторинг
            pkill -f "contract-monitor"
            
            # Удаляем директорию
            rm -rf "$OLD_SITE"
            echo "Старый сайт успешно удален"
        fi
    fi
}

# Функция для проверки готовности нового сайта
check_new_site_files() {
    local missing_files=0
    
    echo "=== Проверка файлов нового сайта ==="
    
    if [ ! -f "index.html" ]; then
        echo "❌ Отсутствует index.html"
        missing_files=1
    else
        echo "✅ index.html найден"
    fi
    
    if [ ! -f "styles.css" ]; then
        echo "❌ Отсутствует styles.css"
        missing_files=1
    else
        echo "✅ styles.css найден"
    fi
    
    if [ ! -f "script.js" ]; then
        echo "❌ Отсутствует script.js"
        missing_files=1
    else
        echo "✅ script.js найден"
    fi
    
    # Проверяем наличие изображений
    if [ ! -d "images" ] && [ ! -f "*.png" ]; then
        echo "⚠️ Предупреждение: не найдены изображения"
    else
        echo "✅ Изображения найдены"
    fi
    
    if [ $missing_files -eq 1 ]; then
        echo "❌ Ошибка: отсутствуют необходимые файлы"
        exit 1
    fi
}

# Основное меню
case "$1" in
    "new")
        check_new_site_files
        create_new_site
        ;;
    "delete")
        delete_old_site
        ;;
    "new-after-delete")
        check_new_site_files
        delete_old_site
        create_new_site
        ;;
    "check")
        check_new_site_files
        ;;
    *)
        echo "Использование:"
        echo "  ./site-manager.sh check         - проверить наличие всех необходимых файлов"
        echo "  ./site-manager.sh new           - создать новый сайт"
        echo "  ./site-manager.sh delete        - удалить текущий сайт"
        echo "  ./site-manager.sh new-after-delete - удалить старый и создать новый сайт"
        ;;
esac 