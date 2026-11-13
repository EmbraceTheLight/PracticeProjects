#!/bin/bash

set -euo pipefail

if (( EUID == 0 )); then
    echo "请以要配置 Oh My Zsh 的普通用户运行此脚本，去掉 sudo。" >&2
    exit 1
fi

account=$(id -un)

retry_git_clone() {
    local url=$1
    local dest=$2
    local count=1
    local max_attempts=3

    # 重试安装时复用之前已克隆成功的插件。
    if [[ -e "$dest" || -L "$dest" ]]; then
        if [[ -d "$dest/.git" ]] && git -C "$dest" rev-parse --verify HEAD >/dev/null; then
            echo "插件已存在，跳过克隆：$dest"
            return 0
        fi
        echo "插件目录已存在，但不是完整的 Git 仓库：$dest" >&2
        return 1
    fi

    while [ $count -le $max_attempts ]; do
        echo "尝试第 $count 次克隆 $url ..."
        git clone "$url" "$dest" && return 0

        echo "克隆失败，尝试使用镜像..."
        git clone "https://ghfast.top/$url" "$dest" && return 0

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
fi

if [[ ! -x /bin/zsh ]]; then
    echo "找不到可执行的 /bin/zsh，无法完成安装。" >&2
    exit 1
fi

# 安装 oh-my-zsh
# 使用 export 确保子进程能正确继承环境变量
# RUNZSH=no  防止安装完后自动启动 zsh 导致后续命令不执行
# CHSH=no    等插件和主题配置完成后，再切换登录 Shell
if [[ ! -d "$HOME/.oh-my-zsh" ]]; then
    installer=$(mktemp)
    trap 'rm -f -- "$installer"' EXIT
    curl -fsSL https://gitee.com/pocmon/ohmyzsh/raw/master/tools/install.sh -o "$installer"
    if [[ ! -s "$installer" ]]; then
        echo "Oh My Zsh 安装脚本下载结果为空。" >&2
        exit 1
    fi
    RUNZSH=no CHSH=no sh "$installer"
    rm -f -- "$installer"
    trap - EXIT
fi

# 校验 oh-my-zsh 安装是否成功
if [ ! -f "$HOME/.zshrc" ]; then
    echo "❌ oh-my-zsh 安装失败：~/.zshrc 未生成，退出脚本"
    exit 1
fi

if [ ! -d "$HOME/.oh-my-zsh" ]; then
    echo "❌ oh-my-zsh 安装失败：~/.oh-my-zsh 目录不存在，退出脚本"
    exit 1
fi

echo "✅ oh-my-zsh 安装成功"

# 安装插件
retry_git_clone https://github.com/zsh-users/zsh-autosuggestions "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}"/plugins/zsh-autosuggestions
retry_git_clone https://github.com/zsh-users/zsh-syntax-highlighting.git "${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}"/plugins/zsh-syntax-highlighting

# 修改 plugins 配置
echo "🔧 配置 plugins..."
sed -i 's/^plugins=.*/plugins=(git zsh-autosuggestions zsh-syntax-highlighting z sudo colored-man-pages)/' ~/.zshrc

# 配置 oh-my-zsh 主题
echo "🔧 配置主题为 itchy..."
sed -i 's/^ZSH_THEME=.*/ZSH_THEME="itchy"/' ~/.zshrc

# 设置 gnome-terminal 配置
#echo "🔧 导入 gnome-terminal 配置..."
#dconf load /org/gnome/terminal/ < "$SCRIPT_DIR/gnome-terminal.conf"

account_entry=$(getent passwd "$account")
if [[ "${account_entry##*:}" != /bin/zsh ]]; then
    echo "🔄 设置 zsh 为登录 Shell..."
    chsh -s /bin/zsh
    account_entry=$(getent passwd "$account")
    if [[ "${account_entry##*:}" != /bin/zsh ]]; then
        echo "登录 Shell 未切换到 /bin/zsh。" >&2
        exit 1
    fi
fi

echo ""
echo "✅ 全部配置完成！正在切换到 zsh..."

# 用 exec 替换当前 shell 进程为 zsh，直接进入配置好的环境
exec zsh -l
