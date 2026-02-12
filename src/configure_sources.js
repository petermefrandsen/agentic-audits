const fs = require('fs');

// Allow passing config path as argument or env var
const configPath = process.argv[2] || process.env.SOURCES_CONFIG;

if (!configPath) {
    // No config provided, return empty defaults
    console.log(JSON.stringify({ mcpServers: {}, mcpPackages: [], webSources: '' }));
    process.exit(0);
}

try {
    const content = fs.readFileSync(configPath, 'utf8');
    const lines = content.split('\n');
    const sources = [];
    let current = null;

    for (const line of lines) {
        const nameMatch = line.match(/^\s*-\s*name:\s*(.+)/);
        if (nameMatch) {
            if (current) sources.push(current);
            current = { name: nameMatch[1].trim() };
            continue;
        }
        if (!current) continue;
        const typeMatch = line.match(/^\s*type:\s*(.+)/);
        const packageMatch = line.match(/^\s*package:\s*["']?([^"']+)["']?/);
        const urlMatch = line.match(/^\s*url:\s*["']?([^"']+)["']?/);
        const enabledMatch = line.match(/^\s*enabled:\s*(.+)/);
        if (typeMatch) current.type = typeMatch[1].trim();
        if (packageMatch) current.package = packageMatch[1].trim();
        if (urlMatch) current.url = urlMatch[1].trim();
        if (enabledMatch) current.enabled = enabledMatch[1].trim() === 'true';
    }
    if (current) sources.push(current);

    // Build MCP servers object from enabled mcp sources
    const mcpServers = {};
    const mcpPackages = [];
    for (const s of sources) {
        if (s.type === 'mcp' && s.enabled && s.package) {
            mcpServers[s.name] = { command: 'npx', args: ['-y', s.package] };
            mcpPackages.push(s.package);
        }
    }

    // Collect enabled web source URLs
    const webUrls = sources
        .filter(s => s.type === 'web' && s.enabled && s.url)
        .map(s => s.url);

    // Output as JSON for the shell to consume
    const output = {
        mcpServers: mcpServers,
        mcpPackages: mcpPackages,
        webSources: webUrls.length > 0
            ? 'Also consult these documentation sources: ' + webUrls.join(', ')
            : ''
    };
    console.log(JSON.stringify(output));
} catch (e) {
    // On error (e.g. file not found), return empty defaults
    // We log to stderr for debugging but stdout must remain valid JSON
    console.error(`Error parsing sources config: ${e.message}`);
    console.log(JSON.stringify({ mcpServers: {}, mcpPackages: [], webSources: '' }));
}
