#include "defs.h"
#include <ctype.h>
#include <string.h>

#ifndef NULL
#define NULL	((void *) 0)
#endif

int nroff = 1;

#define NROFF (-666)
#define TROFF (-667)

STRDEF *chardef, *strdef;
LONGSTRDEF *defdef;
INTDEF *intdef;

static INTDEF standardint[] = {
    { V('n',' '), NROFF, 0, NULL },
    { V('t',' '), TROFF, 0, NULL },
    { V('o',' '), 1,     0, NULL },
    { V('e',' '), 0,     0, NULL },
    { V('.','l'), 70,    0, NULL },
    { V('.','$'), 0,     0, NULL },
    { V('.','A'), NROFF, 0, NULL },
    { V('.','T'), TROFF, 0, NULL },
    { V('.','V'), 1,     0, NULL }, /* the me package tests for this */
    { 0, 0, 0, NULL } };

static STRDEF standardstring[] = {
    { V('<','='), 2, "&lt;=", NULL  }, /* less equal */
    { V('>','='), 2, "&gt=;", NULL  }, /* greather equal */
    { V('A','m'), 1, "&amp;", NULL  }, /* infinity */
    { V('B','a'), 1, "|", NULL  }, /* vartical bar */
    { V('G','e'), 2, "&gt=;", NULL  }, /* greather equal */
    { V('G','t'), 1, "&gt;", NULL  }, /* greather than */
    { V('I','f'), 1, "&infin;", NULL  }, /* infinity */
    { V('L','e'), 2, "&lt;=", NULL  }, /* less equal */
    { V('L','q'), 1, "&ldquo;", NULL  }, /* left double quote  */
    { V('L','t'), 1, "&lt;", NULL  }, /* less than */
    { V('N','a'), 3, "NaN", NULL  }, /* not a number */
    { V('N','e'), 2, "!=", NULL  }, /* not equal */
    { V('P','i'), 2, "Pi", NULL  }, /* pi */
    { V('P','m'), 1, "&plusmn;", NULL  }, /* plus minus */
    { V('R',' '), 1, "&#174;", NULL },
    { V('R','q'), 1, "&rdquo;", NULL  }, /* right double quote  */
    { V('a','a'), 1, "'", NULL  }, /* accute accent  */
    { V('g','a'), 1, "`", NULL  }, /* grave accent  */
    { V('l','q'), 2, "``", NULL },
    { V('q',' '), 1, "&quot;", NULL  }, /* straight double quote  */
    { V('r','q'), 2, "''", NULL },
    { V('u','a'), 1, "^", NULL  }, /* upwards arrow  */
    { 0, 0, NULL, NULL}
};

static STRDEF standardchar[] = {
    { V('*','*'), 1, "*", NULL  },	/* math star */
    { V('*','A'), 1, "&Alpha;", NULL },
    { V('*','B'), 1, "&Beta;", NULL },
    { V('*','C'), 1, "&Xi;", NULL },
    { V('*','D'), 1, "&Delta;", NULL },
    { V('*','E'), 1, "&Epsilon;", NULL },
    { V('*','F'), 1, "&Phi;", NULL },
    { V('*','G'), 1, "&Gamma;", NULL },
    { V('*','H'), 1, "&Theta;", NULL },
    { V('*','I'), 1, "&Iota;", NULL },
    { V('*','K'), 1, "&Kappa;", NULL },
    { V('*','L'), 1, "&Lambda;", NULL },
    { V('*','M'), 1, "&Mu;", NULL },
    { V('*','N'), 1, "&Nu;", NULL },
    { V('*','O'), 1, "&Omicron;", NULL },
    { V('*','P'), 1, "&Pi;", NULL },
    { V('*','Q'), 1, "&Psi;", NULL },
    { V('*','R'), 1, "&Rho;", NULL },
    { V('*','S'), 1, "&Sigma;", NULL },
    { V('*','T'), 1, "&Tau;", NULL },
    { V('*','U'), 1, "&Upsilon;", NULL },
    { V('*','W'), 1, "&Omega;", NULL },
    { V('*','X'), 1, "&Chi;", NULL },
    { V('*','Y'), 1, "&Eta;", NULL },
    { V('*','Z'), 1, "&Zeta;", NULL },
    { V('*','a'), 1, "&alpha;", NULL },
    { V('*','b'), 1, "&beta;", NULL },
    { V('*','c'), 1, "&xi;", NULL },
    { V('*','d'), 1, "&delta;", NULL },
    { V('*','e'), 1, "&epsilon;", NULL },
    { V('*','f'), 1, "&phi;", NULL },
    { V('*','g'), 1, "&gamma;", NULL },
    { V('*','h'), 1, "&theta;", NULL },
    { V('*','i'), 1, "&iota;", NULL },
    { V('*','k'), 1, "&kappa;", NULL },
    { V('*','l'), 1, "&lambda;", NULL },
    { V('*','m'), 1, "&mu;", NULL },
    { V('*','n'), 1, "&nu;", NULL },
    { V('*','o'), 1, "&omicron;", NULL },
    { V('*','p'), 1, "&pi;", NULL },
    { V('*','q'), 1, "&psi;", NULL },
    { V('*','r'), 1, "&rho;", NULL },
    { V('*','s'), 1, "&sigma;", NULL },
    { V('*','t'), 1, "&tau;", NULL },
    { V('*','u'), 1, "&upsilon;", NULL },
    { V('*','w'), 1, "&omega;", NULL },
    { V('*','x'), 1, "&chi;", NULL },
    { V('*','y'), 1, "&eta;", NULL },
    { V('*','z'), 1, "&zeta;", NULL },
    { V('\'','A'), 1, "&Aacute;", NULL },
    { V('\'','E'), 1, "&Eacute;", NULL },
    { V('\'','I'), 1, "&Iacute;", NULL },
    { V('\'','O'), 1, "&Oacute;", NULL },
    { V('\'','U'), 1, "&Uacute;", NULL },
    { V('\'','Y'), 1, "&Yacute;", NULL },
    { V('\'','a'), 1, "&aacute;", NULL },
    { V('\'','e'), 1, "&eacute;", NULL },
    { V('\'','i'), 1, "&iacute;", NULL },
    { V('\'','o'), 1, "&oacute;", NULL },
    { V('\'','u'), 1, "&uacute;", NULL },
    { V('\'','y'), 1, "&yacute;", NULL },
    { V('!','='), 1, "&ne;", NULL },
    { V('%','0'), 1, "&permil;", NULL },
    { V('+','-'), 1, "&plusmn;", NULL },
    { V(',','C'), 1, "&Ccedil;", NULL },
    { V(',','c'), 1, "&ccedil;", NULL },
    { V('-','>'), 1, "&rarr;", NULL },
    { V('-','D'), 1, "&ETH;", NULL },
    { V('.','i'), 1, "&#x131;", NULL },
    { V('/','L'), 1, "&#x141;", NULL },
    { V('/','O'), 1, "&Oslash;", NULL },
    { V('/','l'), 1, "&#x142;", NULL },
    { V('/','o'), 1, "&oslash;", NULL },
    { V('1','2'), 1, "&#189;", NULL  },
    { V('1','4'), 1, "&#188;", NULL  },
    { V('3','4'), 1, "&#190;", NULL  },
    { V(':','A'), 1, "&Auml;", NULL },
    { V(':','E'), 1, "&Euml;", NULL },
    { V(':','I'), 1, "&Iuml;", NULL },
    { V(':','O'), 1, "&Ouml;", NULL },
    { V(':','U'), 1, "&Uuml;", NULL },
    { V(':','a'), 1, "&auml;", NULL },
    { V(':','e'), 1, "&euml;", NULL },
    { V(':','i'), 1, "&iuml;", NULL },
    { V(':','o'), 1, "&ouml;", NULL },
    { V(':','u'), 1, "&uuml;", NULL },
    { V(':','y'), 1, "&yuml;", NULL },
    { V('<','-'), 1, "&larr;", NULL },
    { V('<','='), 1, "&le;", NULL },
    { V('<','>'), 1, "&harr;", NULL },
    { V('=','='), 1, "&equiv;", NULL },
    { V('=','~'), 1, "&cong;", NULL },
    { V('>','='), 1, "&ge;", NULL },
    { V('A','E'), 1, "&AElig;", NULL },
    { V('A','h'), 1, "&alepfsym;", NULL },
    { V('C','R'), 1, "&#x240d;", NULL },
    { V('C','s'), 1, "&curren;", NULL },
    { V('D','o'), 1, "$", NULL },
    { V('E','u'), 1, "&euro;", NULL },
    { V('F','c'), 1, "&raquo;", NULL  },
    { V('F','i'), 3, "ffi", NULL  },
    { V('F','l'), 3, "ffl", NULL  },
    { V('F','o'), 1, "&laquo;", NULL  },
    { V('O','E'), 1, "&OElig;", NULL },
    { V('P','o'), 1, "&pound;", NULL },
    { V('S','1'), 1, "&sup1;", NULL },
    { V('S','2'), 1, "&sup2;", NULL },
    { V('S','3'), 1, "&sup3;", NULL },
    { V('S','d'), 1, "&eth;", NULL },
    { V('T','P'), 1, "&THORN;", NULL },
    { V('T','p'), 1, "&thorn;", NULL },
    { V('Y','e'), 1, "&yen;", NULL },
    { V('^','A'), 1, "&Acirc;", NULL },
    { V('^','E'), 1, "&Ecirc;", NULL },
    { V('^','I'), 1, "&Icirc;", NULL },
    { V('^','O'), 1, "&Ocirc;", NULL },
    { V('^','U'), 1, "&Ucirc;", NULL },
    { V('^','a'), 1, "&acirc;", NULL },
    { V('^','e'), 1, "&ecirc;", NULL },
    { V('^','i'), 1, "&icirc;", NULL },
    { V('^','o'), 1, "&ocirc;", NULL },
    { V('^','u'), 1, "&ucirc;", NULL },
    { V('`','A'), 1, "&Agrave;", NULL },
    { V('`','E'), 1, "&Egrave;", NULL },
    { V('`','I'), 1, "&Igrave;", NULL },
    { V('`','O'), 1, "&Ograve;", NULL },
    { V('`','U'), 1, "&Ugrave;", NULL },
    { V('`','a'), 1, "&agrave;", NULL },
    { V('`','e'), 1, "&egrave;", NULL },
    { V('`','i'), 1, "&igrave;", NULL },
    { V('`','o'), 1, "&ograve;", NULL },
    { V('`','u'), 1, "&ugrave;", NULL },
    { V('a','a'), 1, "&acute;", NULL },
    { V('a','e'), 1, "&aelig;", NULL },
    { V('a','p'), 1, "&asymp;", NULL },
    { V('a','q'), 1, "'", NULL },
    { V('a','t'), 1, "@", NULL },
    { V('a','~'), 1, "~", NULL },
    { V('b','a'), 1, "|", NULL },
    { V('b','b'), 1, "|", NULL },
    { V('b','r'), 1, "|", NULL  },
    { V('b','r'), 1, "|", NULL },
    { V('b','u'), 1, "&bull;", NULL },
    { V('b','v'), 1, "|", NULL  },
    { V('c','*'), 1, "&otimes;", NULL },
    { V('c','+'), 1, "&oplus;", NULL },
    { V('c','i'), 1, "&#x25cb;", NULL },
    { V('c','o'), 1, "&#169;", NULL  },
    { V('c','q'), 1, "'", NULL },
    { V('c','t'), 1, "&#162;", NULL  },
    { V('d','A'), 1, "&dArr;", NULL },
    { V('d','a'), 1, "&darr;", NULL },
    { V('d','d'), 1, "=", NULL },
    { V('d','e'), 1, "&#176;", NULL  },
    { V('d','g'), 1, "-", NULL },
    { V('d','i'), 1, "&#247;", NULL  },
    { V('d','q'), 1, "&quot;", NULL  },
    { V('e','m'), 3, "---", NULL  }, 	/* em dash */
    { V('e','n'), 1, "-", NULL }, 	/* en dash */
    { V('e','q'), 1, "=", NULL },
    { V('e','s'), 1, "&#216;", NULL  },
    { V('e','u'), 1, "&euro;", NULL },
    { V('f','/'), 1, "&frasl;", NULL },
    { V('f','c'), 1, "&rsaquo;", NULL  },
    { V('f','f'), 2, "ff", NULL  },
    { V('f','i'), 2, "fi", NULL  },
    { V('f','l'), 2, "fl", NULL  },
    { V('f','m'), 1, "&#180;", NULL  },
    { V('f','o'), 1, "&lsaquo;", NULL  },
    { V('g','a'), 1, "`", NULL  },
    { V('h','A'), 1, "&hArr;", NULL },
    { V('h','y'), 1, "-", NULL  },
    { V('i','f'), 1, "&infin;", NULL },
    { V('i','s'), 8, "Integral", NULL }, /* integral sign */
    { V('l','A'), 1, "&lArr;", NULL },
    { V('l','B'), 1, "[", NULL },
    { V('l','C'), 1, "{", NULL },
    { V('l','a'), 1, "&lt;", NULL },
    { V('l','b'), 1, "[", NULL  },
    { V('l','c'), 2, "|&#175;", NULL  },
    { V('l','f'), 2, "|_", NULL  },
    { V('l','h'), 1, "&#x261a;", NULL },
    { V('l','k'), 1, "<FONT SIZE=\"+2\">{</FONT>", NULL  },
    { V('l','q'), 1, "\"", NULL },
    { V('l','z'), 1, "&loz;", NULL },
    { V('m','c'), 1, "&micro;", NULL },
    { V('m','i'), 1, "-", NULL  },
    { V('m','u'), 1, "&#215;", NULL  },
    { V('n','o'), 1, "&#172;", NULL  },
    { V('o','A'), 1, "&Aring;", NULL },
    { V('o','a'), 1, "&aring;", NULL },
    { V('o','e'), 1, "&oelig;", NULL },
    { V('o','q'), 1, "'", NULL },
    { V('o','r'), 1, "|", NULL },
    { V('p','d'), 1, "d", NULL }, 	/* partial derivative */
    { V('p','l'), 1, "+", NULL },
    { V('p','s'), 1, "&para;", NULL },
    { V('r','!'), 1, "&iexcl;", NULL },
    { V('r','?'), 1, "&iquest;", NULL },
    { V('r','A'), 1, "&rArr;", NULL },
    { V('r','B'), 1, "]", NULL },
    { V('r','C'), 1, "}", NULL },
    { V('r','a'), 1, "&gt;", NULL },
    { V('r','c'), 2, "&#175;|", NULL  },
    { V('r','f'), 2, "_|", NULL  },
    { V('r','g'), 1, "&#174;", NULL  },
    { V('r','h'), 1, "&#x261b;", NULL },
    { V('r','k'), 1, "<FONT SIZE=\"+2\">}</FONT>", NULL  },
    { V('r','n'), 1, "&#175;", NULL  },
    { V('r','q'), 1, "\"", NULL },
    { V('r','s'), 1, "\\", NULL },
    { V('r','u'), 1, "_", NULL },
    { V('s','c'), 1, "&#167;", NULL  },
    { V('s','h'), 1, "#", NULL },
    { V('s','l'), 1, "/", NULL },
    { V('s','q'), 1, "&#x25a1;", NULL },
    { V('s','s'), 1, "&szlig;", NULL },
    { V('t','f'), 1, "&there4;", NULL },
    { V('t','i'), 1, "~", NULL },
    { V('t','m'), 1, "&trade;", NULL },
    { V('t','s'), 1, "s", NULL }, 	/* should be terminal sigma */
    { V('u','A'), 1, "&uArr;", NULL },
    { V('u','a'), 1, "&uarr;", NULL },
    { V('u','l'), 1, "_", NULL },
    { V('~','A'), 1, "&Atilde;", NULL },
    { V('~','N'), 1, "&Ntilde;", NULL },
    { V('~','O'), 1, "&Otilde;", NULL },
    { V('~','a'), 1, "&atilde;", NULL },
    { V('~','n'), 1, "&ntilde;", NULL },
    { V('~','o'), 1, "&otilde;", NULL },
    { 0, 0, NULL, NULL  }

    
};

void stdinit(void) {
    STRDEF *stdf;
    int i;

    stdf = &standardchar[0];
    i = 0;
    while (stdf->nr) {
	if (stdf->st) stdf->st = xstrdup(stdf->st);
	stdf->next = &standardchar[i];
	stdf = stdf->next;
	i++;
    }
    chardef=&standardchar[0];

    stdf=&standardstring[0];
    i=0;
    while (stdf->nr) {
	 /* waste a little memory, and make a copy, to avoid
	    the segfault when we free non-malloced memory */
	if (stdf->st) stdf->st = xstrdup(stdf->st);
	stdf->next = &standardstring[i];
	stdf = stdf->next;
	i++;
    }
    strdef=&standardstring[0];

    intdef=&standardint[0];
    i=0;
    while (intdef->nr) {
	if (intdef->nr == NROFF) intdef->nr = nroff; else
	if (intdef->nr == TROFF) intdef->nr = !nroff;
	intdef->next = &standardint[i];
	intdef = intdef->next;
	i++;
    }
    intdef = &standardint[0];
    defdef = NULL;
}


LONGSTRDEF* find_longstrdef(LONGSTRDEF * head, int nr, char * longname, char ** out_longname)
{
	char *p, c;
	LONGSTRDEF *de;
	
	p = longname;
	while (p && !isspace(*p)) p++;
	c = *p;
	*p = 0;

	de = head;
	while (de && (de->nr != nr || (de->longname && strcmp(longname, de->longname))))
		de = de->next;

	if (out_longname)
		*out_longname = de ? de->longname : xstrdup(longname);
	*p = c;
	return de;
}
