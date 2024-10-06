import type { Challs } from "./types.ts";

export function buildMdTable(_challs: Challs): string {
  let md = "| ID | Name | Genre | Score | Solver |\n";
  md += "|---|---|---|---|---|\n";
  const challs = Object.entries(_challs).sort((a, b) => {
    if (a[1].genre < b[1].genre) {
      return -1;
    }
    if (a[1].genre > b[1].genre) {
      return 1;
    }
    return 0;
  });
  for (const [id, chall] of challs) {
    md +=
      `| ${id} | ${chall.name} | ${chall.genre} | ${chall.score} | ${chall.solver} |\n`;
  }
  return md;
}

export function buildMdSections(challs: Challs): string {
  let md = "";
  const genres = Array.from(
    new Set(Object.values(challs).map((chall) => chall.genre)),
  );
  for (const genre of genres) {
    md += `## ${genre}\n\n`;
    const genreChalls = Object.entries(challs).filter((chall) =>
      chall[1].genre === genre
    );
    for (const [_, chall] of genreChalls) {
      md += `### ${chall.name} (${chall.score}pt / ${chall.solver} solves)\n\n`;
    }
  }
  return md;
}

export function buiildJson(challs: Challs): string {
  return JSON.stringify(challs);
}
