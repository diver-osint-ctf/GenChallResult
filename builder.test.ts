import { buiildJson, buildMdSections, buildMdTable } from "./builder.ts";
import type { Challs } from "./types.ts";
import { assertEquals } from "jsr:@std/assert";

const emptyChalls: Challs = {};

const singleGenreChalls: Challs = {
  c1: { name: "A", genre: "Web", score: 100, solver: 15 },
  c2: { name: "B", genre: "Web", score: 200, solver: 3 },
  c3: { name: "C", genre: "Web", score: 150, solver: 8 },
};

const sameSolverChalls: Challs = {
  c1: { name: "A", genre: "Crypto", score: 100, solver: 5 },
  c2: { name: "B", genre: "Crypto", score: 200, solver: 5 },
};

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
      "| chall3 | Challenge 3 | Crypto | 150 | 8 |\n" +
      "| chall1 | Challenge 1 | Crypto | 100 | 10 |\n" +
      "| chall2 | Challenge 2 | Web | 200 | 5 |\n";
    assertEquals(result, expected);
  },
);

Deno.test(
  "buildMdSections - should generate markdown sections grouped by genre",
  () => {
    const result = buildMdSections(sampleChalls);
    const expected = "## Crypto\n\n" +
      "### Challenge 3 (150pt / 8 solves)\n\n" +
      "### Challenge 1 (100pt / 10 solves)\n\n" +
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

// 空入力テスト
Deno.test("buildMdTable - should return header only for empty input", () => {
  const result = buildMdTable(emptyChalls);
  const expected = "| ID | Name | Genre | Score | Solver |\n" +
    "|---|---|---|---|---|\n";
  assertEquals(result, expected);
});

Deno.test(
  "buildMdSections - should return empty string for empty input",
  () => {
    assertEquals(buildMdSections(emptyChalls), "");
  },
);

Deno.test("buiildJson - should return empty object for empty input", () => {
  assertEquals(buiildJson(emptyChalls), "{}");
});

// 単一ジャンル内でのsolver昇順ソートテスト
Deno.test(
  "buildMdTable - should sort by solver ascending within single genre",
  () => {
    const result = buildMdTable(singleGenreChalls);
    const expected = "| ID | Name | Genre | Score | Solver |\n" +
      "|---|---|---|---|---|\n" +
      "| c2 | B | Web | 200 | 3 |\n" +
      "| c3 | C | Web | 150 | 8 |\n" +
      "| c1 | A | Web | 100 | 15 |\n";
    assertEquals(result, expected);
  },
);

Deno.test(
  "buildMdSections - should sort by solver ascending within single genre",
  () => {
    const result = buildMdSections(singleGenreChalls);
    const expected = "## Web\n\n" +
      "### B (200pt / 3 solves)\n\n" +
      "### C (150pt / 8 solves)\n\n" +
      "### A (100pt / 15 solves)\n\n";
    assertEquals(result, expected);
  },
);

// 同一solver数テスト
Deno.test(
  "buildMdTable - should handle same solver count within genre",
  () => {
    const result = buildMdTable(sameSolverChalls);
    const dataLines = result.split("\n").filter((l) => l.startsWith("| c"));
    assertEquals(dataLines.length, 2);
    assertEquals(dataLines.some((l) => l.includes("| A |")), true);
    assertEquals(dataLines.some((l) => l.includes("| B |")), true);
  },
);
