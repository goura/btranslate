# btranslate
CLI to call Microsoft Translator (Bing Translator) HTTP API

# How to build

    % go build btranslate.go


# How to run

You have to [sign up for Microsoft Translator and get Client ID and Client Secret for you](http://blogs.msdn.com/b/translation/p/gettingstarted1.aspx).

    % export BTRANSLATE_CLIENT_ID="..."
    % export BTRANSLATE_CLIENT_SECRET="..."
    % btranslate -text="愛はさだめ、さだめは死" -from=ja -to=en
    Love is destiny, fate is death.

    % btranslate -text="愛はさだめ、さだめは死" -from=ja -to=en -round_trip
    愛は運命、運命は死。

    % btranslate -text="愛はさだめ、さだめは死" -from=ja -to=en -json
    {"from":"ja","original":"愛はさだめ、さだめは死","round_tripped":"愛は運命、運命は死。","to":"en","translated":"Love is destiny, fate is death."}
