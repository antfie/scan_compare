# Veracode Scan Compare

This is an unofficial Veracode product. It does not come with any support or warrenty.

Use this tool to compare two Veracode scans.

## Usage

We recommend you configure a Veracode API credentials file as documented here: https://docs.veracode.com/r/c_configure_api_cred_file.

You can also use environment variables (`VERACODE_API_KEY_ID` and `VERACODE_API_KEY_SECRET`) or CLI flags (`-vid` and `-vkey`) to authenticate with the Veracode APIs.

```
./scan_compare -h
Scan Compare v1.0
Copyright © Veracode, Inc. 2022. All Rights Reserved.
This is an unofficial Veracode product. It does not come with any support or warrenty.

Usage of scan_compare:
  -a string
        Veracode Platform URL for scan 'A'
  -b string
        Veracode Platform URL for scan 'B'
  -vid string
        Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds
  -vkey string
        Veracode API key - See https://docs.veracode.com/r/t_create_api_creds
```

## Example Usage

```
./scan_compare -a https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview:75603:793744:22132159:22103486:22119136::::5000002 -b https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview:75603:793744:22131974:22103301:22118951::::4999988
```

If you know the build ids you can use them instead of the URL like so:

```
./scan_compare -a 22132159 -b 22131974
```

## Example Output

```
./scan_compare -a 22132159 -b 22131974
Scan Compare v1.0
Copyright © Veracode, Inc. 2022. All Rights Reserved.
This is an unofficial Veracode product. It does not come with any support or warrenty.

Comparing scan 'A' (Build id = 22132159) against scan 'B' (Build id = 22131974)

In common with both scans
=========================
Application: Secure File Transfer
Files uploaded: 0
Top-level modules selected for analysis: 1

Scan 'A'
========
Sandbox: 'Production'
Scan Name: '15 Nov 2022 Static'
Review Modules URL: https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:75603:793744:22132159:22103486:22119136::::5000002
Files uploaded: 127
Submitted: 2022-11-15 12:28:29 +0000 GMT
Duration: 9m4s
Flaws: 49 total, 49 not mitigated

Scan 'B'
========
Sandbox: 'Development'
Scan Name: '15 Nov 2022 Static'
Review Modules URL: https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:75603:793744:22131974:22103301:22118951::::4999988
Files uploaded: 1
Submitted: 2022-11-15 12:29:03 +0000 GMT
Duration: 3m15s
Flaws: 13 total, 13 not mitigated

Differences of Top-Level Modules Selected As An Entry Point
===========================================================
Only in A: 'app.war' - Size = 2MB, Issues = 0, MD5 = 803c155a360d219460d2a59c81389833, Compiler = JAVAC_11, OS = Java J2SE 11, Architecture = JVM
Only in B: 'app-new.war' - Size = 1KB, Issues = 0, MD5 = da0099a578876c08b473a0df8ec589f4, Compiler = JAVAC_11, OS = Java J2SE 11, Architecture = JVM
```

## Development

### Running

```
go run *.go
```

### Compiling

./release.sh