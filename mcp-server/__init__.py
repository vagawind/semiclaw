#!/usr/bin/env python3
"""
SemiClaw MCP Server Package

A Model Context Protocol server that provides access to the SemiClaw knowledge management API.
"""

__version__ = "1.0.0"
__author__ = "SemiClaw Team"
__description__ = "SemiClaw MCP Server - Model Context Protocol server for SemiClaw API"

from .semiclaw_mcp_server import SemiClawClient, run

__all__ = ["SemiClawClient", "run"]
