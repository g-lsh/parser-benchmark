#!/usr/bin/env node

const { parseHTML } = require('linkedom');
const {readFile} = require("node:fs").promises;

async function parsePage() {
    const arg = process.argv[2];

    if (!arg) {
        console.error('❌ arg is required');
        process.exit(1);
    }

    try {
        let body;
        if (arg.startsWith("https://")) {
            const res = await fetch(arg);
            body = await res.text();
            if (!res.ok) {
                throw new Error(`Failed to fetch page: ${res.statusText}`);
            }
        } else {
            body = await readFile(arg, "utf-8");
        }


        const startTime = performance.now();
        const { document } = parseHTML(body);
        document.title = 'HELLO THERE!'
        document.toString()
        const endTime = performance.now();

        console.log('✅ Parsing and writing Time:', (endTime - startTime).toFixed(2), 'ms');

    } catch (error) {
        console.error('❌ Error:', error.message);
    }
}

parsePage();
