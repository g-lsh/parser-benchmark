#!/usr/bin/env node

import { HTMLRewriter } from "html-rewriter-wasm";
import { performance } from "perf_hooks";
import fs from "fs/promises";


const encoder = new TextEncoder();
const decoder = new TextDecoder();

let output = "";
const rewriter = new HTMLRewriter((outputChunk) => {
    output += decoder.decode(outputChunk);
});

rewriter.on("title", {
    element(element) {
        element.setInnerContent("new title");
    },
});

const arg = process.argv[2];

if (!arg) {
    console.error("Please provide a URL or a local file path as an argument.");
    process.exit(1);
}

try {
    let body;
    if (arg.startsWith("https://")) {
        const res = await fetch(arg);
        body = await res.text();
    } else {
        body = await fs.readFile(arg, "utf-8");
    }
    const start = performance.now();
    await rewriter.write(encoder.encode(body));
    await rewriter.end();
    const end = performance.now();
    console.log(`Execution time: ${end - start} milliseconds`);
} finally {
    rewriter.free(); // Remember to free memory
}