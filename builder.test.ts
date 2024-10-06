import { buiildJson, buildMdSections, buildMdTable } from "./builder.ts";
import type { Challs } from "./types.ts";
import { assertEquals } from "jsr:@std/assert";

const sampleChalls: Challs = {
  chall1: { name: "Challenge 1", genre: "Crypto", score: 100, solver: 10 },
  chall2: { name: "Challenge 2", genre: "Web", score: 200, solver: 5 },
  chall3: { name: "Challenge 3", genre: "Crypto", score: 150, solver: 8 },
};

Deno.test(
  "buildMdTable - should generate a markdown table sorted by genre",
  () => {
    const result = buildMdTable(sampleChalls);
    const expected = "| ID | Name | Genre | Score | Solver |\n" +
      "|---|---|---|---|---|\n" +
      "| chall1 | Challenge 1 | Crypto | 100 | 10 |\n" +
      "| chall3 | Challenge 3 | Crypto | 150 | 8 |\n" +
      "| chall2 | Challenge 2 | Web | 200 | 5 |\n";
    assertEquals(result, expected);
  },
);

Deno.test(
  "buildMdSections - should generate markdown sections grouped by genre",
  () => {
    const result = buildMdSections(sampleChalls);
    const expected = "## Crypto\n\n" +
      "### Challenge 1 (100pt / 10 solves)\n\n" +
      "### Challenge 3 (150pt / 8 solves)\n\n" +
      "## Web\n\n" +
      "### Challenge 2 (200pt / 5 solves)\n\n";
    assertEquals(result, expected);
  },
);

Deno.test(
  "buiildJson - should generate a JSON string of the challenges",
  () => {
    const result = buiildJson(sampleChalls);
    const expected = JSON.stringify(sampleChalls);
    assertEquals(result, expected);
  },
);
