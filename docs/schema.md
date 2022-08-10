# Schema

## Why
Maintaining a lot of database columns is boring at best, problematic at worst.
Also exporting / importing settings of game server and making sure that everything is set up correctly and validated is troublesome.
One thing that works great in k8s is schema of resource. Let's make similar in exp too.

## Plans and pricing

### First steps
This is very powerful feature for users.
If we give this to users for free and then take it back or put it behind paywall
it can make significant damage to our rep
So far just hide it from users and think more about this

### Plan for later
Probably offering raw read/write access to schema of game server in free plan would be too much.
Main reason is supporting custom apps/game servers by writing custom commands to it.

Let's keep it at least in lowest paid plan.
That would ensure that everyone that is somehow earning money on our product and receiving not our support is paying us.

For free users it would be enough to use our panel which will give selective access to schema fields
by our panel validation and frontend.

Preview of paid feature: let's just show schema without ability to modify it.
This would allow to export servers easily and import them in other accounts
- users could create "game packages" with free access game server types at even free plan
- but it would require at least the cheapest plan to import
  - we want easy creation of basic game servers so this is kinda problematic
    - importing in free plan is like giving editing for free
