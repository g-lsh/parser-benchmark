#!/usr/bin/env node

const htmlparser2 = require('htmlparser2');
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
        const handler = new htmlparser2.DomHandler();
        const parser = new htmlparser2.Parser(handler);
        parser.write(body);
        parser.end();

        const titleElement = handler.dom.find(el => el.name === 'title');
        if (titleElement && titleElement.children[0]) {
            titleElement.children[0].data = 'hello world';
        }
        const endTime = performance.now();
        console.log('✅ Parsing and writing Time:', (endTime - startTime).toFixed(2), 'ms');

    } catch (error) {
        console.error('❌ Error:', error.message);
    }
}

parsePage();
