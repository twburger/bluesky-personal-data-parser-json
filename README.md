# BlueSky CAR File Data Parser
This project provides a tool to parse and extract data from CAR (Content Addressable Archive) files exported from repositories on the AT Protocol, including platforms like Bluesky. The tool decodes the repository data, processes it into JSON format, and organizes it by lexicon type.


## What is a CAR File?
<img src="https://github.com/user-attachments/assets/e47c4868-d99f-4b5a-bf48-5a970a77292a" width="350" align="right">

A **CAR (Content Addressable Archive)** file is a snapshot of a repository’s state in the AT Protocol. It encapsulates all public data for a repository, such as posts, likes, and social graphs. CAR files are useful for:
- **Portability**: Sharing and offline analysis of repository data.
- **Consistency**: Capturing the exact state of a repository at a specific time.
- **Backup and Archiving**: Preserving repository data independently of live APIs.

<br />

### Summary of Use Cases
| **Use Case**               | **CAR File** | **DID/API Queries** |
|----------------------------|--------------|---------------------|
| Offline access             | ✅           | ❌                  |
| Efficient batch processing | ✅           | ❌                  |
| Real-time updates          | ❌           | ✅                  |
| Backup and archival        | ✅           | ❌                  |
| Selective data retrieval   | ❌           | ✅                  |
| Consistent snapshot of data| ✅           | ❌                  |


## Features

- Decodes CAR files and extracts records in CBOR and JSON formats.
- Aggregates records by lexicon types (e.g., `app.bsky.feed.post`).
- Outputs structured JSON files for each lexicon type.
- Supports offline processing of repository snapshots.

## Requirements

To use this parser, ensure you have the following installed:

- **Go** (1.20 or later)
- A CAR file to process (exported from a Bluesky repository).

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/thomasafink/bluesky-personal-data-parser-json.git
   cd bluesky-personal-data-parser-json
   ```

2. Install the required Go modules:
   ```bash
   go mod tidy
   ```

## Usage

1. Download at https://bsky.app/settings via "Export My Data" and place the CAR file you want to parse in the root directory of the project. Make sure it is named`repo.car` or update the `carFilePath` in the `main.go` file.

<img width="1084" alt="Screenshot 2024-11-19 at 00 52 55" src="https://github.com/user-attachments/assets/5427d178-3621-4ff8-9634-37556c593a28">

2. Run the parser:
   ```bash
   go run main.go
   ```

3. The tool will:
   - Save individual CBOR and JSON records in a directory named after the repository DID.
   - Generate aggregated JSON files for each lexicon type (e.g., `app_bsky_feed_post.json`) in the root directory.

### Example Output

After running the tool, you will find:
- A directory named after the DID (e.g., `did:plc:abc123`) containing:
  - CBOR files for each record.
  - JSON files for each record.
- Aggregated JSON files in the root directory for each lexicon type.

#### Example Aggregated JSON File (`app_bsky_feed_post.json`)
```json
[
  {
    "$type": "app.bsky.feed.post",
    "content": "Hello, world!",
    "createdAt": "2023-07-01T23:30:08.840Z",
    "author": "did:plc:xyz456"
  }
]
```

## Known Issues

- This tool only processes public CAR files. Private or encrypted repositories are not supported.
- Lexicon types not conforming to expected formats may generate warnings.

## Contributing

Feel free to submit issues or contribute improvements via pull requests. Contributions are welcome!

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
