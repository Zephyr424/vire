// Vire v0.1 Alpha - 修正版（关键字优先）
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <windows.h>

#define MAX_VARS 256
long long vars[MAX_VARS] = {0};

char *src;
int pos;
int len;

void error(const char *msg) {
    printf("错误: %s\n", msg);
    exit(1);
}

void skip_ws() {
    while (pos < len && isspace(src[pos])) pos++;
}

long long parse_num() {
    long long val = 0;
    while (pos < len && isdigit(src[pos])) {
        val = val * 10 + (src[pos] - '0');
        pos++;
    }
    return val;
}

char parse_var() {
    if (pos < len && isalpha(src[pos])) {
        return src[pos++];
    }
    return 0;
}

long long parse_expr() {
    skip_ws();
    long long left;
    if (isdigit(src[pos])) {
        left = parse_num();
    } else if (isalpha(src[pos])) {
        char var = parse_var();
        left = vars[(unsigned char)var];
    } else {
        error("期望数字或变量");
    }
    skip_ws();
    while (pos < len && (src[pos] == '+' || src[pos] == '-' || src[pos] == '*' || src[pos] == '/')) {
        char op = src[pos++];
        skip_ws();
        long long right;
        if (isdigit(src[pos])) {
            right = parse_num();
        } else if (isalpha(src[pos])) {
            char var = parse_var();
            right = vars[(unsigned char)var];
        } else {
            error("期望数字或变量");
        }
        skip_ws();
        switch (op) {
            case '+': left += right; break;
            case '-': left -= right; break;
            case '*': left *= right; break;
            case '/': left /= right; break;
        }
    }
    return left;
}

void parse_statement() {
    skip_ws();
    if (pos >= len) return;

    // 1. 处理 print(...) —— 必须优先！
    if (strncmp(&src[pos], "print", 5) == 0) {
        pos += 5;
        skip_ws();
        if (src[pos] == '(') {
            pos++; // 跳过 '('
            long long val = parse_expr();
            skip_ws();
            if (src[pos] == ')') pos++;
            printf("%lld\n", val);
        }
        return;
    }

    // 2. 处理 syscall(...) —— 必须优先！
    if (strncmp(&src[pos], "syscall", 7) == 0) {
        pos += 7;
        skip_ws();
        if (src[pos] == '(') {
            pos++; // 跳过 '('
            long long num = parse_expr();
            skip_ws();
            // 简化跳过其余参数（逗号和数字）
            while (pos < len && src[pos] != ')') pos++;
            if (src[pos] == ')') pos++;
            if (num == 39) {
                printf("%d\n", GetCurrentProcessId());
            } else {
                printf("0\n");
            }
        }
        return;
    }

    // 3. 处理赋值: a = 表达式
    if (isalpha(src[pos])) {
        char var_name = src[pos++];
        skip_ws();
        if (src[pos] == '=') {
            pos++; // 跳过 '='
            long long val = parse_expr();
            vars[(unsigned char)var_name] = val;
        } else {
            error("无效的赋值语句");
        }
        return;
    }

    // 无法识别，跳过
    pos++;
}

int main(int argc, char** argv) {
    if (argc < 2) {
        printf("Vire v0.1 Alpha\n");
        printf("用法: vire.exe 脚本.vire\n");
        return 1;
    }

    FILE* f = fopen(argv[1], "r");
    if (!f) {
        printf("无法打开文件: %s\n", argv[1]);
        return 1;
    }

    fseek(f, 0, SEEK_END);
    long size = ftell(f);
    rewind(f);

    src = (char*)malloc(size + 1);
    if (!src) { printf("内存分配失败\n"); return 1; }

    fread(src, 1, size, f);
    src[size] = '\0';
    fclose(f);

    len = size;
    pos = 0;

    while (pos < len) {
        parse_statement();
        while (pos < len && (src[pos] == '\n' || src[pos] == ';' || src[pos] == ' ' || src[pos] == '\r')) {
            pos++;
        }
    }

    free(src);
    return 0;
}