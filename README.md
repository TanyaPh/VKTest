Когда завершите задачу, в этом README опишите свой ход мыслей: как вы пришли к решению, какие были варианты и почему выбрали именно этот. 

# Что нужно сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение. 

# Необязательно, но было бы круто

Хорошо, если добавите поддержку конфигурации итоговой реализации. Параметры — на ваше усмотрение.

# Xод мыслей

Первое что пришло в голову создать структуру, в которой храниться map с ключем user ID и значением количество проведенных проверок. Эта задумка не подходила, потому что при запуске нескольких программ данные для одних и техже пользователей не суммировались. Необходимо общее хранилище.

Следущая мысль сорханять количество проверок в базе и использовать таймер на заданное время. Этот вариант неподходит тем что считает количество запросов за N секунд и требует перезапуск, а не последнии N секунд. 

Поэтому приняла решение в базе хранить время проверки и считать количество записей за последние N секунд. 