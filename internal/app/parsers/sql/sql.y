%{
package sql
%}

%union {
    ival string
    str string
}

%type <ival> Iconst

%token <ival> ICONST
%token <str> BCONST

%%

Iconst:
    ICONST
    {
        $$ = $1;
    };