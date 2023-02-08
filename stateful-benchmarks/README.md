## Introduction

The goal of this benchmark is to figure out how much space we can save by not sending symbols each time we send a pprof file.

## Methodology

* This benchmark takes a series of pprof files as input
* For each file, it generates 3 derivative profiles:
  * one with symbols removed
  * one with symbols encoded normally as strings (same as input)
  * one with symbols replaces by a hash of the symbol (currently I use murmur3, take 8 bytes and encode it as a hex string)
* Then it encodes each of the 3 profiles back to protobuf bytes, and compresses it with gzip
* The program then outputs the size of each of the 3 profiles
* In addition to that, it also deduplicates all symbols from all profiles, puts it in a separate pprof file, and compresses it with gzip, and then prints the size of that as well

## Preliminary Findings

I ran this benchmark on a few different applications, and the results are as follows:
* Most apps had profiles of 10KB - 50KB size
* Symbols typically take 30-60% of the total size of the pprof file
* For 50-100 profiles, deduplicated symbols typically take 10KB - 50KB of space, which is pretty good, it means there's a lot of duplication in the symbols
* In theory, if we find a way to not send symbols each time, we'd get 30-60% reduction in size of each payload
* In practice however, my naive implementation of hashing the symbols did not improve anything, the numbers we're getting are pretty close to the numbers we get when we send symbols normally
* My understanding is that the reason for this is that gzip compression is effective on symbols, and not effective on hashes (because entropy is already high there)


## Next Steps

Naive hashing of symbols doesn't seem to work. I think the next thing we should explore is hashing of the entire stacktraces. We could modify this benchmark to replace each stacktrace with a hash of the stacktrace, then aggregate all unique stacktraces similarly to how we're currently aggregating unique strings, and see if that improves the numbers.

### Usage

To run the benchmarks, you need to cd into this directory (`opentelemetry-profiling/stateful-benchmarks`), have go installed, and then run the following command:

```
go run . <path to your pprof files>
# for example
go run . inputs/app1/*.pb.gz
```

### Example Outputs

Go app with no labels, low CPU utilization:
```
noSymbols       symbolsAsStrings        symbolsAsStrings percent difference     hashedSymbols     hashedSymbols percent difference
3170    5930    +87.07% 5754    +81.51%
3072    6007    +95.54% 5708    +85.81%
7915    13282   +67.81% 13270   +67.66%
9572    15781   +64.87% 15926   +66.38%
3334    6633    +98.95% 6239    +87.13%
2454    4638    +89.00% 4536    +84.84%
9977    16589   +66.27% 16745   +67.84%
2460    4487    +82.40% 4435    +80.28%
8449    14184   +67.88% 14175   +67.77%
2191    4193    +91.37% 4022    +83.57%
2698    5099    +88.99% 4986    +84.80%
3164    5804    +83.44% 5769    +82.33%
10178   16797   +65.03% 16824   +65.30%
9692    16150   +66.63% 16236   +67.52%
9522    15771   +65.63% 15922   +67.21%
10216   17152   +67.89% 17097   +67.36%
10265   16634   +62.05% 16828   +63.94%
2761    5483    +98.59% 5197    +88.23%
2263    4326    +91.16% 4229    +86.88%
2047    3958    +93.36% 3809    +86.08%
3925    7515    +91.46% 7174    +82.78%

separateSymbols: 15292
```

Go app with many labels, high CPU utilization (~70%):
```
noSymbols       symbolsAsStrings        symbolsAsStrings percent difference     hashedSymbols     hashedSymbols percent difference
23073   33221   +43.98% 32628   +41.41%
22656   32726   +44.45% 31928   +40.93%
29630   41437   +39.85% 40644   +37.17%
26618   38990   +46.48% 37984   +42.70%
30779   42707   +38.75% 41890   +36.10%
24546   35673   +45.33% 34946   +42.37%
23211   34129   +47.04% 33273   +43.35%
21143   30229   +42.97% 29449   +39.28%
23297   34047   +46.14% 33320   +43.02%
... truncated for brevity
21802   31485   +44.41% 30563   +40.18%
20766   29740   +43.21% 28904   +39.19%
21512   31854   +48.08% 30974   +43.98%
29219   41476   +41.95% 40433   +38.38%
28996   41365   +42.66% 40236   +38.76%
26935   38389   +42.52% 37447   +39.03%
26145   37753   +44.40% 36833   +40.88%
26551   37237   +40.25% 36532   +37.59%
22453   33366   +48.60% 32364   +44.14%
20017   29435   +47.05% 28348   +41.62%

separateSymbols: 33486
```

Another go app with many labels, high CPU utilization (~70%):
```
noSymbols       symbolsAsStrings        symbolsAsStrings percent difference     hashedSymbols     hashedSymbols percent difference
32774   45489   +38.80% 44303   +35.18%
33929   46827   +38.01% 45805   +35.00%
47810   62403   +30.52% 61220   +28.05%
48932   64812   +32.45% 63579   +29.93%
57788   72838   +26.04% 71514   +23.75%
29339   42034   +43.27% 40953   +39.59%
34244   47169   +37.74% 46109   +34.65%
11774   19885   +68.89% 18970   +61.12%
54919   70444   +28.27% 69222   +26.04%
... truncated for brevity
30914   43260   +39.94% 42229   +36.60%
50649   65874   +30.06% 64667   +27.68%
44879   57682   +28.53% 57071   +27.17%
43436   57996   +33.52% 56892   +30.98%
52412   67778   +29.32% 66573   +27.02%
40936   53524   +30.75% 52572   +28.42%
33498   46326   +38.29% 45230   +35.02%
28461   40580   +42.58% 39502   +38.79%

separateSymbols: 42386
```

Go app with no labels, high CPU utilization (~70%):
```
noSymbols       symbolsAsStrings        symbolsAsStrings percent difference     hashedSymbols     hashedSymbols percent difference
29669   43872   +47.87% 42744   +44.07%
25377   38429   +51.43% 37235   +46.73%
27789   40953   +47.37% 40106   +44.32%
27783   40868   +47.10% 40024   +44.06%
10046   16936   +68.58% 16872   +67.95%
24707   37605   +52.20% 36490   +47.69%
22227   33222   +49.47% 32599   +46.66%
25432   37943   +49.19% 37032   +45.61%
24145   36374   +50.65% 35384   +46.55%
24782   37274   +50.41% 36339   +46.63%
26002   38518   +48.13% 37784   +45.31%
7948    13855   +74.32% 13716   +72.57%
26681   39965   +49.79% 38876   +45.71%
30650   45033   +46.93% 43852   +43.07%

separateSymbols: 25747
```

