#!/bin/bash

retry_git_clone() {
    local url=$1
    local dest=$2
    local count=1
    local max_attempts=3

    while [ $count -le $max_attempts ]; do
        echo "尝试第 $count 次克隆 $url ..."
        git clone "$url" "$dest" && return 0

        echo "克隆失败，尝试使用镜像..."
        git clone "https://ghproxy.com/$url" "$dest" && return 0

        echo "等待 5 秒后重试..."
        sleep 5
        count=$((count + 1))
    done

    echo "git clone 多次失败，退出脚本"
    exit 1
}

sudo apt update 

if ! command -v git >/dev/null 2>&1; then
    echo "git 未安装，开始安装..."
    sudo apt install -y git
fi

if ! command -v curl >/dev/null 2>&1; then
    echo "curl 未安装，开始安装..."
    sudo apt install -y curl
fi

if ! command -v zsh >/dev/null 2>&1; then
    echo "zsh 未安装，开始安装..."
    sudo apt install -y zsh

    # 配置 zsh 为默认终端 
    chsh -s /bin/zsh
fi

# 安装 oh-my-zsh
sh -c "$(curl -fsSL https://gitee.com/pocmon/ohmyzsh/raw/master/tools/install.sh)"

# 安装插件
retry_git_clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions 10
retry_git_clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting 10

# 修改 plugins 配置
sed -i 's/^plugins=.*/plugins=(git zsh-autosuggestions zsh-syntax-highlighting z sudo colored-man-pages)/' ~/.zshrc

# 配置 oh-my-zsh 主题
sed -i 's/^ZSH_THEME=.*/ZSH_THEME="itchy"/' ~/.zshrc

# 设置 gnome-terminal 配置
dconf load /org/gnome/terminal/ < gnome-terminal.conf

# 重新加载配置
source ~/.zshrc
