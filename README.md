# Veracode Scan Compare

This is an unofficial Veracode product. It does not come with any support or warrenty.

Use this tool to compare two Veracode scans.

## Usage

We recommend you configure a Veracode API credentials file as documented here: https://docs.veracode.com/r/c_configure_api_cred_file.

Alternatively you can use environment variables (`VERACODE_API_KEY_ID` and `VERACODE_API_KEY_SECRET`) or CLI flags (`-vid` and `-vkey`) to authenticate with the Veracode APIs.

```
./scan_compare -h
Scan Compare v1.0
Copyright © Veracode, Inc. 2022. All Rights Reserved.
This is an unofficial Veracode product. It does not come with any support or warrenty.

Usage of scan_compare:
  -a string
        Veracode Platform URL for scan "A"
  -b string
        Veracode Platform URL for scan "B"
  -vid string
        Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds
  -vkey string
        Veracode API key - See https://docs.veracode.com/r/t_create_api_creds
```

## Example Usage

```
./scan_compare -a https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview:75603:793744:22132159:22103486:22119136::::5000002 -b https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview:75603:793744:22131974:22103301:22118951::::4999988
```

If you know the build ids you can use them instead of URLs like so:

```
./scan_compare -a 22132159 -b 22131974
```

## Example Output

```
./scan_compare -a 22132159 -b 22131974
Scan Compare v1.0
Copyright © Veracode, Inc. 2022. All Rights Reserved.
This is an unofficial Veracode product. It does not come with any support or warrenty.

Comparing scan "A" (Build id = 22132159) against scan "B" (Build id = 22131974)

In common with both scans
=========================
Application: "Secure File Transfer"
Scan name: "15 Nov 2022 Static"
Top-level modules selected for analysis: 1
Engine version: 20221021161836

Scan A
========
Sandbox: "Production"
Review Modules URL: https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:75603:793744:22132159:22103486:22119136::::5000002
Files uploaded: 127
Submitted: 2022-11-15 12:28:29 +0000 GMT
Duration: 9m4s
Flaws: 49 total, 49 not mitigated

Scan B
========
Sandbox: "Development"
Review Modules URL: https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:75603:793744:22131974:22103301:22118951::::4999988
Files uploaded: 1
Submitted: 2022-11-15 12:29:03 +0000 GMT
Duration: 3m15s
Flaws: 13 total, 13 not mitigated

Differences of Top-Level Modules Selected As An Entry Point For Scanning
========================================================================
Only in A: "app.war" - Size = 2MB, Missing Supporting Files = 85, MD5 = 803c155a360d219460d2a59c81389833, Platform = JVM / Java J2SE 11 / JAVAC_11
Only in B: "app-new.war" - Size = 0KB, Missing Supporting Files = 427, MD5 = da0099a578876c08b473a0df8ec589f4, Platform = JVM / Java J2SE 11 / JAVAC_11

Differences of Top-Level Modules Not Selected As An Entry Point (And Not Scanned)
=================================================================================
Only in A: "lib-0.1.0-SNAPSHOT.jar" - Size = 12MB, Missing Supporting Files = 20, MD5 = e81932dcfb6a7bc3df152e5cf6538a6f, Platform = JVM / Java J2SE 11 / JAVAC_11

Summary
========
B was submitted 34s after A
A took longer by 5m49s
```

## Development

### Running

```
go run *.go
```

### Compiling

./release.sh