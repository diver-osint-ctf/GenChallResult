export type Chall = {
  name: string;
  genre: string;
  score: number;
  solver: number;
};

export type Challs = {
  // key: chall id
  [key: string]: Chall;
};
