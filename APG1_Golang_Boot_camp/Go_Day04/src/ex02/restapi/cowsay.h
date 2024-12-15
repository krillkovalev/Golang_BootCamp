#ifndef ASK_COW_H
#define ASK_COW_H

// Подключаем необходимые стандартные библиотеки
#include <stdlib.h>
#include <string.h>

// Объявление функции ask_cow
// Принимает: указатель на строку (массив символов)
// Возвращает: указатель на строку (динамически выделенная память)
char *ask_cow(char phrase[]);

#endif // ASK_COW_H
