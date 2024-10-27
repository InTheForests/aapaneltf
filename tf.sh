#!/bin/bash

if ! command -v curl &> /dev/null; then
    echo "curl 未安装 尝试安装 curl..."
    if [[ -f /etc/debian_version ]]; then
        sudo apt update && sudo apt install -y curl
    elif [[ -f /etc/redhat-release ]]; then
        sudo yum install -y curl
    else
        echo "不支持的系统 请手动安装 curl!"
        exit 1
    fi
fi

if ! command -v tar &> /dev/null; then
    echo "tar 未安装 尝试安装 tar..."
    if [[ -f /etc/debian_version ]]; then
        sudo apt update && sudo apt install -y tar
    elif [[ -f /etc/redhat-release ]]; then
        sudo yum install -y tar
    else
        echo "不支持的系统 请手动安装 tar!"
        exit 1
    fi
fi

os_type=$(uname -s)
architecture=$(uname -m)
latest_url=$(curl -s https://api.github.com/repos/InTheForests/aapaneltf/releases/latest | \
             grep "browser_download_url.*${os_type}_${architecture}.tar.gz" | \
             cut -d '"' -f 4)

if [[ -z "$latest_url" ]]; then
    echo "未找到适用于 ${os_type} ${architecture} 版本下载链接"
    exit 1
fi

echo "正在下载 ${latest_url} ..."
curl -L -o aapaneltf.tar.gz "$latest_url"

echo "下载完成 解压"
tar -xzvf aapaneltf.tar.gz aapaneltf
echo "解压完成 删除 aapaneltf.tar.gz"
rm aapaneltf.tar.gz
echo "运行 aapaneltf"
chmod +x aapaneltf
./aapaneltf
echo "删除 aapaneltf"
rm aapaneltf