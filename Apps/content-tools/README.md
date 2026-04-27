# Google Analytics 4 CLI (ga-cli)

`ga-cli` is a command-line interface tool for querying Google Analytics 4 (GA4) data. It allows you to quickly fetch reports on acquisitions, page views, events, and more directly from your terminal.

## Features

- Query GA4 data using simple commands.
- Support for various report types:
  - Acquisition (Source / Medium)
  - Page views
  - Events
  - Page views by Acquisition Source
  - Page views by Country
- Customizable date ranges.
- Uses Google Application Default Credentials (ADC) for secure authentication.

## Prerequisites

1.  **Google Cloud Project**: You need a Google Cloud project with the [Google Analytics Data API](https://console.cloud.google.com/apis/library/analyticsdata.googleapis.com) enabled.
2.  **Authentication**: Set up [Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/provide-credentials-adc). The easiest way is to use the Google Cloud CLI:
    ```bash
    gcloud auth application-default login
    ```
3.  **GA4 Property ID**: You need the Property ID of your Google Analytics 4 property. You can find this in the GA4 Admin panel under **Property Settings > Property Details**.

## Installation

To install `ga-cli`, you need to have Go installed on your system.

```bash
go install .
```

Alternatively, you can build the binary:

```bash
go build -o ga-cli main.go
```

## Configuration

The tool requires your GA4 Property ID to be set as an environment variable:

```bash
export GA_PROPERTY_ID="your-property-id"
```

## Usage

The main command is `ga-cli search`, which has several subcommands.

### Global Flags

These flags are available for all `search` subcommands:

- `--start-date`: Start date for the report (e.g., `2023-01-01`, `30daysAgo`, `today`). Default: `30daysAgo`.
- `--end-date`: End date for the report (e.g., `2023-01-31`, `today`). Default: `today`.

### Subcommands

#### 1. Acquisition
Shows user acquisition data grouped by Source / Medium.
```bash
ga-cli search acquisition
```

#### 2. Pages
Shows page views for different page paths.
```bash
ga-cli search pages
```

#### 3. Events
Shows event counts for different event names.
```bash
ga-cli search events
```

#### 4. Pages by Acquisition
Shows page views broken down by page path and session source/medium.
```bash
ga-cli search pages-by-acquisition
```

#### 5. Pages by Country
Shows page views broken down by page path and country.
```bash
ga-cli search pages-by-country
```

#### 6. Version
Check the version of the tool.
```bash
ga-cli version
```

## Example

Get page views for the last 7 days:

```bash
ga-cli search pages --start-date=7daysAgo --end-date=today
```

## License

[MIT License](LICENSE) (if applicable)
