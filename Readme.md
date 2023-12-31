# Introduction & Explanation

BetAlgorithm is a simple cli-based arbitrage bet detection and return calculation algorithm. The algorithm uses free data from [The Odds API](https://the-odds-api.com/) to determine if there exists opportunities to place bets that have guaranteed return.

I do not warrant that this works as it is based purely on [theory](http://www.aussportsbetting.com/guide/sports-betting-arbitrage/) and I have not used it in practice, nor do I condone or recommend that you do either. Gambling is risky and no one should do it. This is simply a code-based demonstration of the arbitrage betting formula.

## Installation

```
1. git clone https://github.com/KrachkovskyiYevhen/BetAlgorithm.git
```

2. Create environment variable "THEODDSAPIKEY" with your [The Odds API](https://the-odds-api.com/) Key

```

## Usage

```

1. cd BetAlgorithm
2. go run .

```

## Arguments

- Datatype
- Default
- Descriptions

```

--demo

```
- boolean
- false
- Denotes weather to use the demo data provided/your own or to use the odds api.


```

--demo_file

```
- string
- ./test_data.json
- Specifies the demo data file location relative to the directory. Note: Must be supplied if demo is true


```

--bet

```
- integer
- 1000
- The total amount of money you're willing to wager


```

--sport

```
- string
- upcoming
- The sport type to download from the odds api. Please reference their site for more inputs.


```

--region

```
- string
- au
- The region to download from the odds api. Please reference their site for more inputs.


```

--verbose

```
- boolean
- false
- Whether to log all program output.
```
