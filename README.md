Проект упращающий создания swagger документации 

Запустить проект: make run

Запуск миграций - смотрите в Makefile

для выполнения миграций поставить пакет 
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
Решение проблемы
1. Убедитесь, что вы установили migrate
Выполните команду установки:

go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
Эта команда скачает и установит CLI migrate в директорию $GOPATH/bin.

2. Проверьте наличие бинарного файла
Убедитесь, что файл migrate установлен:

ls $GOPATH/bin/migrate
Если переменная GOPATH не настроена, проверьте стандартный путь для Go:

ls ~/go/bin/migrate
3. Добавьте бинарный файл в $PATH
Если файл migrate найден, добавьте его директорию в переменную окружения $PATH. Например, если файл находится в ~/go/bin:

export PATH=$PATH:~/go/bin
Добавьте эту строку в ваш ~/.zshrc, чтобы изменения сохранились после перезагрузки терминала:

echo 'export PATH=$PATH:~/go/bin' >> ~/.zshrc
source ~/.zshrc
4. Проверьте доступность команды
После настройки $PATH убедитесь, что команда migrate работает:

migrate --help
Если проблемы сохраняются
Если всё ещё возникает ошибка, убедитесь, что Go правильно установлен, и переменная GOPATH настроена:

echo $GOPATH
По умолчанию GOPATH обычно указывает на ~/go. Если это не так, настройте его:

export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
Добавьте настройки в ваш ~/.zshrc:

echo 'export GOPATH=~/go' >> ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc