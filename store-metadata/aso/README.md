# KuranAI ASO source of truth

`store-listing.json` contains the approved Turkish, English, and Arabic store
metadata for App Store Connect and Google Play Console.

The metadata intentionally:

- puts the highest-intent search terms in the app name/title and subtitle;
- keeps Apple keyword fields under 100 characters without spaces after commas;
- avoids duplicated Apple title/subtitle terms in the keyword fields;
- uses relevant phrases naturally instead of keyword blocks;
- describes only features that are present in the released apps;
- avoids ranking, pricing, or unverifiable religious claims;
- keeps the education/research and non-fatwa disclosure visible.

## Measurement cadence

Review search impressions, product-page/store-listing conversion, installs,
ratings, and retention every 14 days. Change only one high-impact variable at a
time in Product Page Optimization or Store Listing Experiments so results remain
interpretable.
