#!/bin/bash

# 卸载 oh-my-zsh 及相关配置的脚本

set -euo pipefail

if (( EUID == 0 )); then
    echo "请以要卸载 Oh My Zsh 的普通用户运行此脚本，去掉 sudo。" >&2
    exit 1
fi

account=$(id -un)

is_zsh_shell() {
    case "${1##*/}" in
        zsh|zsh[0-9]*|rzsh) return 0 ;;
    esac
    [[ -e /bin/zsh && "$1" -ef /bin/zsh ]]
}

echo "⚠️  此脚本将执行以下操作："
echo "  1. 如果当前登录 Shell 是 zsh，先切换回 bash"
echo "  2. 删除 oh-my-zsh (~/.oh-my-zsh)"
echo "  3. 删除 ~/.zshrc（如有备份则恢复）"
echo "  4. 可选：卸载 zsh"
echo ""
read -p "确认继续？(y/N) " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "已取消。"
    exit 0
fi

# 1. 切换登录 Shell，并确认账户记录已更新。
account_entry=$(getent passwd "$account")
if is_zsh_shell "${account_entry##*:}"; then
    if [[ ! -x /bin/bash ]]; then
        echo "找不到可执行的 /bin/bash，已停止卸载。" >&2
        exit 1
    fi
    echo "🔄 将登录 Shell 切换回 /bin/bash ..."
    chsh -s /bin/bash
    account_entry=$(getent passwd "$account")
    if [[ "${account_entry##*:}" != /bin/bash ]]; then
        echo "无法确认登录 Shell 已切换到 /bin/bash，已停止卸载。" >&2
        exit 1
    fi
fi

# 2. 删除 oh-my-zsh
if [ -d "$HOME/.oh-my-zsh" ]; then
    echo "🗑️  删除 ~/.oh-my-zsh ..."
    rm -rf "$HOME/.oh-my-zsh"
else
    echo "ℹ️  ~/.oh-my-zsh 不存在，跳过。"
fi

# 3. 处理 ~/.zshrc
if [ -f "$HOME/.zshrc.pre-oh-my-zsh" ]; then
    echo "♻️  恢复 oh-my-zsh 安装前的 ~/.zshrc 备份..."
    mv "$HOME/.zshrc.pre-oh-my-zsh" "$HOME/.zshrc"
elif [ -f "$HOME/.zshrc" ]; then
    echo "🗑️  删除 ~/.zshrc ..."
    rm -f "$HOME/.zshrc"
else
    echo "ℹ️  ~/.zshrc 不存在，跳过。"
fi

# 4. 可选：卸载 zsh
echo ""
read -p "是否卸载 zsh？(y/N) " uninstall_zsh
if [[ "$uninstall_zsh" == "y" || "$uninstall_zsh" == "Y" ]]; then
    # 分别检查本地账户和 NSS 可枚举账户。NSS 可能关闭某些来源的枚举。
    all_accounts=$(getent passwd)
    if [[ -z "$all_accounts" || ! -r /etc/passwd ]]; then
        echo "无法检查系统账户，已保留 zsh。" >&2
        exit 1
    fi
    all_accounts=$(sort -u /etc/passwd - <<< "$all_accounts")
    zsh_accounts=()
    while IFS=: read -r name _ _ _ _ _ shell; do
        if is_zsh_shell "$shell"; then
            zsh_accounts+=("$name ($shell)")
        fi
    done <<< "$all_accounts"
    if (( ${#zsh_accounts[@]} > 0 )); then
        echo "以下账户仍使用 zsh 作为登录 Shell，已保留 zsh：" >&2
        printf '  %s\n' "${zsh_accounts[@]}" >&2
        exit 1
    fi
    echo "🗑️  卸载 zsh ..."
    sudo apt remove -y zsh
    echo "✅ zsh 已卸载。"
else
    echo "ℹ️  保留 zsh。"
fi

echo ""
echo "✅ 卸载完成！请重新打开终端使更改生效。"

