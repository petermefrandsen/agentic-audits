import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { exec } from 'child_process';
import fs from 'fs';
import path from 'path';
import util from 'util';

const execPromise = util.promisify(exec);
const scriptPath = path.resolve(__dirname, '../configure_sources.js');
const tempConfigFile = path.resolve(__dirname, 'temp_sources.yml');

describe('configure_sources.js', () => {
    afterAll(() => {
        if (fs.existsSync(tempConfigFile)) {
            fs.unlinkSync(tempConfigFile);
        }
    });

    it('should return empty defaults if no config file provided', async () => {
        const { stdout } = await execPromise(`node ${scriptPath}`);
        const result = JSON.parse(stdout);
        expect(result).toEqual({ mcpServers: {}, mcpPackages: [], webSources: '' });
    });

    it('should parse valid YAML config correctly', async () => {
        const yamlContent = `
- name: github-mcp-server
  type: mcp
  package: "@github/mcp-server"
  enabled: true
- name: ignored-server
  type: mcp
  package: "ignore-me"
  enabled: false
- name: docs
  type: web
  url: "https://docs.example.com"
  enabled: true
`;
        fs.writeFileSync(tempConfigFile, yamlContent);

        const { stdout } = await execPromise(`node ${scriptPath} ${tempConfigFile}`);
        const result = JSON.parse(stdout);

        expect(result.mcpServers).toHaveProperty('github-mcp-server');
        expect(result.mcpServers['github-mcp-server']).toEqual({
            command: 'npx',
            args: ['-y', '@github/mcp-server']
        });
        expect(result.mcpPackages).toContain('@github/mcp-server');
        expect(result.mcpServers).not.toHaveProperty('ignored-server');
        expect(result.webSources).toContain('https://docs.example.com');
    });

    it('should handle malformed YAML gracefully', async () => {
        fs.writeFileSync(tempConfigFile, 'invalid yaml content');
        // The script prints error to stderr but acts as if empty config on stdout
        const { stdout } = await execPromise(`node ${scriptPath} ${tempConfigFile}`);
        const result = JSON.parse(stdout);
        expect(result).toEqual({ mcpServers: {}, mcpPackages: [], webSources: '' });
    });
});
