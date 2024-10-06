import { buiildJson, buildMdSections, buildMdTable } from "./builder.ts";
import { loadChalls, loadSolves } from "./loadFiles.ts";

function exportFile(content: string, fileName: string) {
  const encoder = new TextEncoder();
  const data = encoder.encode(content);
  Deno.writeFile(`out/${fileName}`, data);
}

async function main() {
  if (Deno.args.length < 2) {
    console.log(
      "Usage: deno task run <challenges.csv> <solves.csv>",
    );
    console.log(
      "Example: deno task run ~/Download/HogeCTF-challenges.csv ~/Download/HogeCTF-solves.csv",
    );
    Deno.exit(1);
  }
  const challs = await loadChalls(Deno.args[0]);
  const challsWithSolves = await loadSolves(Deno.args[1], challs);

  const mdTable = buildMdTable(challsWithSolves);
  exportFile(mdTable, "summary.md");
  const mdSections = buildMdSections(challsWithSolves);
  exportFile(mdSections, "sections.md");
  const json = buiildJson(challsWithSolves);
  exportFile(json, "challs.json");
}
main();
