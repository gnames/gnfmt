# GNcsv

This library provides reading and writing to CSV, TSV and PSV (pipe-separated)
files.

## Usage

### Reader

```go
import csvConfig "github.com/gnames/gnfmt/gncsv/config"
...
opts := []csvConfig.Option{
    csvConfig.OptPath(path),
    csvConfig.OptBadRowMode(cfg.BadRow),
    csvConfig.OptWithQuotes(cfg.WithQuotes),
}
chIn := make(chan []string)
// create new config with required options.
cfg, err := csvConfig.New(opts...)
if err != nil {
    return err
}

go func() {
    ...
    for row := range chIn {
      ...  
    }
}()

csv := gncsv.New(csvCfg)
_, err = csv.Read(context.Background(), chIn)
if err != nil {
    return err
}
close(chIn)
...
```

### Writer

```go
headers := []string{"id","name"}
opts := []config.Option{
    config.OptPath("write.csv")
    config.OptHeaders("headers")
cfg, err := config.New(opts...)
...
w.gncsv.New(cfg.Write) ch := make(chan []string)
go func() {
        defer wg.Done()
        err = w.WriteStream(context.Background(), ch)
        ...
}

}
```

