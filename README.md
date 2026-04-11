[![test](https://github.com/diver-osint-ctf/GenChallResult/actions/workflows/test.yaml/badge.svg)](https://github.com/diver-osint-ctf/GenChallResult/actions/workflows/test.yaml)

# GenChallResult

Writeupのために、3種類の形式で問題情報を出力するツール

```md
## OSINT

### テスト問題 (500pt / 1 solves)

### テスト問題2 (530pt / 0 solves)

## Web

### テスト問題3 (510pt / 0 solves)
```

```md
| ID | Name        | Genre | Score | Solver |
| -- | ----------- | ----- | ----- | ------ |
| 1  | テスト問題  | OSINT | 500   | 1      |
| 2  | テスト問題2 | OSINT | 530   | 0      |
| 3  | テスト問題3 | Web   | 510   | 0      |
```

```json
{
  "1": { "name": "テスト問題", "genre": "OSINT", "score": 500, "solver": 1 },
  "2": { "name": "テスト問題2", "genre": "OSINT", "score": 530, "solver": 0 },
  "3": { "name": "テスト問題3", "genre": "Web", "score": 510, "solver": 0 }
}
```

## How to use

### Docker Compose (推奨)

```bash
docker compose run --rm app -c challenges.csv -s solves.csv -t teams.csv
```

### Go

```bash
# ビルド
make build

# 実行
./genChallResult -c challenges.csv -s solves.csv -t teams.csv
```

`-t` (teams.csv) はオプションです。指定するとhidden/bannedチームのsolve数を除外します。

challenges.csvとsolves.csvとteams.csvは`CTFd > 管理画面 > Config > Import & Export > Download CSV`より、challengesとsolvesとteamsを選択してダウンロードしてください。

![image](./assets/image.png)

## Development

```bash
make test   # テスト実行
make lint   # go vet によるlint
make build  # バイナリビルド
make clean  # バイナリ削除
```
