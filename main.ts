import { parse } from "@std/flags";
import { buiildJson, buildMdSections, buildMdTable } from "./builder.ts";
import { loadChalls, loadSolves, loadTeams } from "./loadFiles.ts";

function exportFile(content: string, fileName: string) {
  const encoder = new TextEncoder();
  const data = encoder.encode(content);
  Deno.writeFile(`out/${fileName}`, data);
}

async function main() {
  const flags = parse(Deno.args, {
    string: ["c", "s", "t", "challenges", "solves", "teams"],
    alias: { c: "challenges", s: "solves", t: "teams" },
  });

  const challengesPath = flags.c || flags.challenges;
  const solvesPath = flags.s || flags.solves;
  const teamsPath = flags.t || flags.teams;

  if (!challengesPath || !solvesPath) {
    console.log(
      "Usage: deno task run -c <challenges.csv> -s <solves.csv> [-t <teams.csv>]",
    );
    console.log(
      "Example: deno task run -c ~/Download/HogeCTF-challenges.csv -s ~/Download/HogeCTF-solves.csv -t ~/Download/HogeCTF-teams.csv",
    );
    Deno.exit(1);
  }

  const challs = await loadChalls(challengesPath);
  let ignoredTeams = new Set<string>();
  if (teamsPath) {
    ignoredTeams = await loadTeams(teamsPath);
  }
  const challsWithSolves = await loadSolves(solvesPath, challs, ignoredTeams);

  const mdTable = buildMdTable(challsWithSolves);
  exportFile(mdTable, "summary.md");
  const mdSections = buildMdSections(challsWithSolves);
  exportFile(mdSections, "sections.md");
  const json = buiildJson(challsWithSolves);
  exportFile(json, "challs.json");
}
main();
