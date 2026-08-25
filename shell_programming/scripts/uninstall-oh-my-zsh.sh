#!/bin/bash

# 卸载 oh-my-zsh 及相关配置的脚本

set -e

echo "⚠️  此脚本将执行以下操作："
echo "  1. 删除 oh-my-zsh (~/.oh-my-zsh)"
echo "  2. 删除 ~/.zshrc（如有备份则恢复）"
echo "  3. 将默认 shell 切换回 bash"
echo "  4. 重置 gnome-terminal 配置为默认值"
echo "  5. 可选：卸载 zsh"
echo ""
read -p "确认继续？(y/N) " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "已取消。"
    exit 0
fi

# 1. 删除 oh-my-zsh
if [ -d "$HOME/.oh-my-zsh" ]; then
    echo "🗑️  删除 ~/.oh-my-zsh ..."
    rm -rf "$HOME/.oh-my-zsh"
else
    echo "ℹ️  ~/.oh-my-zsh 不存在，跳过。"
fi

# 2. 处理 ~/.zshrc
if [ -f "$HOME/.zshrc.pre-oh-my-zsh" ]; then
    echo "♻️  恢复 oh-my-zsh 安装前的 ~/.zshrc 备份..."
    mv "$HOME/.zshrc.pre-oh-my-zsh" "$HOME/.zshrc"
elif [ -f "$HOME/.zshrc" ]; then
    echo "🗑️  删除 ~/.zshrc ..."
    rm -f "$HOME/.zshrc"
else
    echo "ℹ️  ~/.zshrc 不存在，跳过。"
fi

# 3. 切换默认 shell 回 bash
current_shell="$(getent passwd "$USER" | cut -d: -f7)"
if [[ "$current_shell" == */zsh ]]; then
    echo "🔄 将默认 shell 切换回 /bin/bash ..."
    chsh -s /bin/bash
    echo "✅ 默认 shell 已切换为 bash。"
else
    echo "ℹ️  当前默认 shell 不是 zsh ($current_shell)，跳过。"
fi

# 4. 重置 gnome-terminal 配置
#echo "🔄 重置 gnome-terminal 配置为默认值..."
#dconf reset -f /org/gnome/terminal/legacy/profiles:/

# 5. 可选：卸载 zsh
echo ""
read -p "是否卸载 zsh？(y/N) " uninstall_zsh
if [[ "$uninstall_zsh" == "y" || "$uninstall_zsh" == "Y" ]]; then
    echo "🗑️  卸载 zsh ..."
    sudo apt remove -y zsh
    sudo apt autoremove -y
    echo "✅ zsh 已卸载。"
else
    echo "ℹ️  保留 zsh。"
fi

echo ""
echo "✅ 卸载完成！请重新打开终端使更改生效。"

