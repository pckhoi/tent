

SQL =
  stmts:Stmt+

Stmt =
  Comment* _ stmt:( SetStmt / CreateTableStmt / CreateSeqStmt / CreateExtensionStmt / CreateTypeStmt / AlterTableStmt / AlterSeqStmt / CommentExtensionStmt )


/*
 ██████╗██████╗ ███████╗ █████╗ ████████╗███████╗    ████████╗ █████╗ ██████╗ ██╗     ███████╗
██╔════╝██╔══██╗██╔════╝██╔══██╗╚══██╔══╝██╔════╝    ╚══██╔══╝██╔══██╗██╔══██╗██║     ██╔════╝
██║     ██████╔╝█████╗  ███████║   ██║   █████╗         ██║   ███████║██████╔╝██║     █████╗  
██║     ██╔══██╗██╔══╝  ██╔══██║   ██║   ██╔══╝         ██║   ██╔══██║██╔══██╗██║     ██╔══╝  
╚██████╗██║  ██║███████╗██║  ██║   ██║   ███████╗       ██║   ██║  ██║██████╔╝███████╗███████╗
 ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝       ╚═╝   ╚═╝  ╚═╝╚═════╝ ╚══════╝╚══════╝
*/                                                                                            

CreateTableStmt =
  "CREATE"i _1 "TABLE"i _1 tablename:Ident _ "(" _ fields:( FieldDef ( _ "," _ FieldDef )* ) _ ")" _ ";" EOL*

FieldDef =
  name:Ident _1 dataType:DataType constraint:ColumnConstraint?

ColumnConstraint =
  nameOpt:( _1 "CONSTRAINT"i _1 (StringConst / Ident) )? constraintClauses:( NotNullCls / NullCls / CheckCls )+

NotNullCls =
  _1 "NOT"i _1 "NULL"i

NullCls =
  _1 "NULL"i

CheckCls =
  _1 "CHECK"i _1 expr:WrappedExpr noInherit:( _1 "NO"i _1 "INHERIT"i )?

WrappedExpr =
  "(" Expr+ ")"

Expr =
  WrappedExpr / [^()]+


/*
██████╗  █████╗ ████████╗ █████╗ ████████╗██╗   ██╗██████╗ ███████╗
██╔══██╗██╔══██╗╚══██╔══╝██╔══██╗╚══██╔══╝╚██╗ ██╔╝██╔══██╗██╔════╝
██║  ██║███████║   ██║   ███████║   ██║    ╚████╔╝ ██████╔╝█████╗  
██║  ██║██╔══██║   ██║   ██╔══██║   ██║     ╚██╔╝  ██╔═══╝ ██╔══╝  
██████╔╝██║  ██║   ██║   ██║  ██║   ██║      ██║   ██║     ███████╗
╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝      ╚═╝   ╚═╝     ╚══════╝
*/

DataType =
  t:( TimestampT / TimeT / VarcharT / CharT / BitVarT / BitT / IntT / PgOidT / OtherT / CustomT )

TimestampT =
  "TIMESTAMP"i prec:SecPrecision withTimeZone:( WithTZ / WithoutTZ )?

TimeT =
  "TIME"i prec:SecPrecision withTimeZone:( WithTZ / WithoutTZ )?

SecPrecision =
  ( _1 [0-6])?

WithTZ =
  _1 "WITH"i _1 "TIME"i _1 "ZONE"i

WithoutTZ =
  ( _1 "WITHOUT"i _1 "TIME"i _1 "ZONE"i )?

CharT =
  ( "CHARACTER"i / "CHAR"i ) length:( "(" NonZNumber ")" )?

VarcharT =
  ( ( "CHARACTER"i _1 "VARYING"i ) / "VARCHAR"i ) length:( "(" NonZNumber ")" )?

BitT =
  "BIT"i length:( "(" NonZNumber ")" )?

BitVarT =
  "BIT"i _1 "VARYING"i length:( "(" NonZNumber ")" )?

IntT =
  ( "INTEGER"i / "INT"i )

PgOidT =
  ( "OID"i / "REGPROCEDURE"i / "REGPROC"i / "REGOPERATOR"i / "REGOPER"i / "REGCLASS"i / "REGTYPE"i / "REGROLE"i / "REGNAMESPACE"i / "REGCONFIG"i / "REGDICTIONARY"i )

OtherT =
  ( "DATE"i / "SMALLINT"i / "BIGINT"i / "DECIMAL"i / "NUMERIC"i / "REAL"i / "SMALLSERIAL"i / "SERIAL"i / "BIGSERIAL"i / "BOOLEAN"i / "TEXT"i / "MONEY"i / "BYTEA"i / "POINT"i / "LINE"i / "LSEG"i / "BOX"i / "PATH"i / "POLYGON"i / "CIRCLE"i / "CIDR"i / "INET"i / "MACADDR"i / "UUID"i / "XML"i / "JSONB"i / "JSON"i )

CustomT =
  Ident


/*
 ██████╗██████╗ ███████╗ █████╗ ████████╗███████╗    ███████╗███████╗ ██████╗ 
██╔════╝██╔══██╗██╔════╝██╔══██╗╚══██╔══╝██╔════╝    ██╔════╝██╔════╝██╔═══██╗
██║     ██████╔╝█████╗  ███████║   ██║   █████╗      ███████╗█████╗  ██║   ██║
██║     ██╔══██╗██╔══╝  ██╔══██║   ██║   ██╔══╝      ╚════██║██╔══╝  ██║▄▄ ██║
╚██████╗██║  ██║███████╗██║  ██║   ██║   ███████╗    ███████║███████╗╚██████╔╝
 ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝    ╚══════╝╚══════╝ ╚══▀▀═╝
*/

CreateSeqStmt =
  "CREATE"i _1 "SEQUENCE"i _1 name:Ident verses:CreateSeqVerse* _ ";" EOL*

CreateSeqVerse =
  verse:( IncrementBy / MinValue / NoMinValue / MaxValue / NoMaxValue / Start / Cache / Cycle / OwnedBy )

IncrementBy =
  _1 "INCREMENT"i (_1 "BY"i)? _1 num:NonZNumber

MinValue =
  _1 "MINVALUE"i _1 val:NonZNumber

NoMinValue =
  _1 "NO"i _1 "MINVALUE"i

MaxValue =
  _1 "MAXVALUE"i _1 val:NonZNumber

NoMaxValue =
  _1 "NO"i _1 "MAXVALUE"i

Start =
  _1 "START"i (_1 "WITH"i)? _1 start:NonZNumber

Cache =
  _1 "CACHE"i _1 cache:NonZNumber

Cycle =
  no:(_1 "NO"i)? _1 "CYCLE"i

OwnedBy =
  _1 "OWNED"i _1 "BY"i _1 name:( "NONE"i / TableDotCol )

/*
 ██████╗██████╗ ███████╗ █████╗ ████████╗███████╗    ████████╗██╗   ██╗██████╗ ███████╗
██╔════╝██╔══██╗██╔════╝██╔══██╗╚══██╔══╝██╔════╝    ╚══██╔══╝╚██╗ ██╔╝██╔══██╗██╔════╝
██║     ██████╔╝█████╗  ███████║   ██║   █████╗         ██║    ╚████╔╝ ██████╔╝█████╗  
██║     ██╔══██╗██╔══╝  ██╔══██║   ██║   ██╔══╝         ██║     ╚██╔╝  ██╔═══╝ ██╔══╝  
╚██████╗██║  ██║███████╗██║  ██║   ██║   ███████╗       ██║      ██║   ██║     ███████╗
 ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝       ╚═╝      ╚═╝   ╚═╝     ╚══════╝
 */

CreateTypeStmt =
  "CREATE"i _1 "TYPE"i _1 typename:Ident _1 "AS"i _1 typedef:EnumDef _ ";" EOL*

EnumDef =
  "ENUM" _ "(" _ vals:( StringConst ( _ ',' _ StringConst )*  ) _ ")"


 /*
 █████╗ ██╗  ████████╗███████╗██████╗     ████████╗ █████╗ ██████╗ ██╗     ███████╗
██╔══██╗██║  ╚══██╔══╝██╔════╝██╔══██╗    ╚══██╔══╝██╔══██╗██╔══██╗██║     ██╔════╝
███████║██║     ██║   █████╗  ██████╔╝       ██║   ███████║██████╔╝██║     █████╗  
██╔══██║██║     ██║   ██╔══╝  ██╔══██╗       ██║   ██╔══██║██╔══██╗██║     ██╔══╝  
██║  ██║███████╗██║   ███████╗██║  ██║       ██║   ██║  ██║██████╔╝███████╗███████╗
╚═╝  ╚═╝╚══════╝╚═╝   ╚══════╝╚═╝  ╚═╝       ╚═╝   ╚═╝  ╚═╝╚═════╝ ╚══════╝╚══════╝
*/

AlterTableStmt =
  "ALTER"i _1 "TABLE"i _1 name:Ident _1 "OWNER"i _1 "TO"i _1 owner:Ident _ ";" EOL*


/*
 █████╗ ██╗  ████████╗███████╗██████╗     ███████╗███████╗ ██████╗ 
██╔══██╗██║  ╚══██╔══╝██╔════╝██╔══██╗    ██╔════╝██╔════╝██╔═══██╗
███████║██║     ██║   █████╗  ██████╔╝    ███████╗█████╗  ██║   ██║
██╔══██║██║     ██║   ██╔══╝  ██╔══██╗    ╚════██║██╔══╝  ██║▄▄ ██║
██║  ██║███████╗██║   ███████╗██║  ██║    ███████║███████╗╚██████╔╝
╚═╝  ╚═╝╚══════╝╚═╝   ╚══════╝╚═╝  ╚═╝    ╚══════╝╚══════╝ ╚══▀▀═╝ 
*/

AlterSeqStmt =
  "ALTER"i _1 "SEQUENCE"i _1 name:Ident _1 "OWNED"i _1 "BY"i _1 owner:TableDotCol _ ";" EOL*

TableDotCol =
  table:Ident "." column:Ident


/*
 ██████╗ ████████╗██╗  ██╗███████╗██████╗     ███████╗████████╗███╗   ███╗████████╗
██╔═══██╗╚══██╔══╝██║  ██║██╔════╝██╔══██╗    ██╔════╝╚══██╔══╝████╗ ████║╚══██╔══╝
██║   ██║   ██║   ███████║█████╗  ██████╔╝    ███████╗   ██║   ██╔████╔██║   ██║   
██║   ██║   ██║   ██╔══██║██╔══╝  ██╔══██╗    ╚════██║   ██║   ██║╚██╔╝██║   ██║   
╚██████╔╝   ██║   ██║  ██║███████╗██║  ██║    ███████║   ██║   ██║ ╚═╝ ██║   ██║   
 ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝    ╚══════╝   ╚═╝   ╚═╝     ╚═╝   ╚═╝   
*/

CommentExtensionStmt =
  "COMMENT"i _1 "ON"i _1 "EXTENSION"i _ extension:Ident _ "IS"i _ comment:StringConst _ ";" EOL*

CreateExtensionStmt =
  "CREATE"i _1 "EXTENSION"i _1 ( "IF"i _1 "NOT"i _1 "EXISTS"i _1 )? extension:Ident _1 "WITH"i _1 "SCHEMA"i _1 schema:Ident _ ";" EOL*

SetStmt =
  "SET"i _ key:Key _ ( "=" / "TO"i ) _ values:CommaSeparatedValues _ ";" EOL*

Key =
  [a-z_]i+


/*
██╗   ██╗ █████╗ ██╗     ██╗   ██╗███████╗███████╗
██║   ██║██╔══██╗██║     ██║   ██║██╔════╝██╔════╝
██║   ██║███████║██║     ██║   ██║█████╗  ███████╗
╚██╗ ██╔╝██╔══██║██║     ██║   ██║██╔══╝  ╚════██║
 ╚████╔╝ ██║  ██║███████╗╚██████╔╝███████╗███████║
  ╚═══╝  ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚══════╝╚══════╝
*/


CommaSeparatedValues =
  vals:( Value ( _ ',' _ Value )* )

Value =
  ( Number / Boolean / StringConst / Ident )

StringConst =
  "'" value:([^'\n] / "''")* "'"

Ident =
  [a-z_]i [a-z_0-9$]i*

Number =
  ( "0" / [1-9][0-9]* )

NonZNumber =
  [1-9][0-9]*

Boolean =
  value:( BooleanTrue / BooleanFalse )

BooleanTrue =
  ( "TRUE" / "'" BooleanTrueString "'" / BooleanTrueString )

BooleanTrueString =
  ( "true" / "yes" / "on" / "t" / "y" )

BooleanFalse =
  ( "FALSE" / "'" BooleanFalseString "'" / BooleanFalseString )

BooleanFalseString =
  ( "false" / "no" / "off" / "f" / "n" )


/*
███╗   ███╗██╗███████╗ ██████╗
████╗ ████║██║██╔════╝██╔════╝
██╔████╔██║██║███████╗██║     
██║╚██╔╝██║██║╚════██║██║     
██║ ╚═╝ ██║██║███████║╚██████╗
╚═╝     ╚═╝╚═╝╚══════╝ ╚═════╝
*/

Comment =
  ( SingleLineComment / MultilineComment )

MultilineComment =
  "/*" .* "*/" EOL

SingleLineComment =
  "--" [^\r\n]* EOL

EOL =
  [ \t]* ("\r\n" / "\n\r" / "\r" / "\n")

_ "whitespace" =
  [ \t\r\n]*

_1 "at least 1 whitespace" =
  [ \t\r\n]+

EOF =
  !.
