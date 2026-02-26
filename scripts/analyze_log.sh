#!/bin/bash
# 接收一个日志文件名作为参数，调用 Claude Code 对日志进行分析，并将结果(仅纯文本)保存到一个报告文件中

# 检查参数
if [ -z "$1" ];then
    echo "Usage: $0 <log_file_path>"
    exit 1
fi

LOG_FILE=$1
REPORT_FILE="${LOG_FILE}.report.txt"
PROMPT="你是一个资深的 SRE 工程师，请深入分析这份日志文件，找出其中的错误 (Error) 和警告 (Warning/Warn)，对它们进行分类和计数，并总结出最可能导致这些问题的根本原因。以 Markdown 格式输出你的报告。"

echo "Analyzing log file: $LOG_FILE..."

# 使用管道将日志文件内容传递给Claude Code的Headless模式
# 使用 --output-format json 来获取结构化结果
# 使用 jq 提取最终的纯文本报告
# 使用 > 将输出重定向到报告文件
#
cat "$LOG_FILE" | claude-glm -p "$PROMPT" --output-format json | jq -r '.result' > "$REPORT_FILE"

# 检查命令是否执行成功
if [ $? -eq 0 ];then
    echo "Analsis complete.Report saved to: $REPORT_FILE"
    echo "--- Report Preview ---"
    head -n 10 "$REPORT_FILE"
else
    echo "Analsis failed."
    # 尝试打印一些错误信息
    cat "$LOG_FILE" | claude-glm -p "$PROMPT" --output-format json > /dev/null
    exit 1
fi
