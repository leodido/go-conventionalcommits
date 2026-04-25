%%{

machine common;

ws = ' ';

nl = [\n];

dash = '-';

lpar = '(';

rpar = ')';

colon = 0x3A;

exclamation = 0x21;

# high_byte widens Ragel's built-in `print` class ([\x20-\x7e], ASCII
# printable only) with the high-byte range [\x80-\xff] so productions
# that capture user-supplied free-form text (trailer values, scope,
# free-form types) accept any non-control byte transparently.
#
# The class is "any high byte", NOT "any valid UTF-8 byte". It admits
# every byte in [\x80-\xff] including bytes that never appear in
# valid UTF-8 (\xc0, \xc1, \xfe, \xff), lone leaders, lone
# continuations, and overlong-encoding leaders. Validating that
# captured slices form well-formed UTF-8 is the caller's
# responsibility; an opt-in `WithStrictUTF8` option for that lives
# in a follow-up PR.
#
# \x7f (DEL) is intentionally excluded; it is a control byte. C0
# controls [\x00-\x1f] (except \n where the production allows it)
# remain rejected.
#
# See https://github.com/leodido/go-conventionalcommits/issues/50
high_byte = print | 0x80..0xff;

}%%