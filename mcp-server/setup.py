#!/usr/bin/env python3
"""
SemiClaw MCP Server 安装脚本
"""

from setuptools import setup


# 读取 README 文件
def read_readme():
    try:
        with open("README.md", "r", encoding="utf-8") as f:
            return f.read()
    except FileNotFoundError:
        return "SemiClaw MCP Server - Model Context Protocol server for SemiClaw API"


# 读取依赖
def read_requirements():
    try:
        with open("requirements.txt", "r", encoding="utf-8") as f:
            return [
                line.strip() for line in f if line.strip() and not line.startswith("#")
            ]
    except FileNotFoundError:
        return ["mcp>=1.0.0", "requests>=2.31.0"]


setup(
    name="semiclaw-mcp-server",
    version="1.0.0",
    author="SemiClaw Team",
    author_email="support@semiclaw.com",
    description="SemiClaw MCP Server - Model Context Protocol server for SemiClaw API",
    long_description=read_readme(),
    long_description_content_type="text/markdown",
    url="https://github.com/vagawind/semiclaw-mcp",
    py_modules=["semiclaw_mcp_server", "main", "run_server", "run", "test_module"],
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Internet :: WWW/HTTP :: HTTP Servers",
        "Topic :: Scientific/Engineering :: Artificial Intelligence",
    ],
    python_requires=">=3.10",
    install_requires=read_requirements(),
    entry_points={
        "console_scripts": [
            "semiclaw-mcp-server=main:sync_main",
            "semiclaw-server=run_server:main",
        ],
    },
    include_package_data=True,
    data_files=[
        ("", ["README.md", "requirements.txt", "LICENSE"]),
    ],
    keywords="mcp model-context-protocol semiclaw knowledge-management api-server",
)
