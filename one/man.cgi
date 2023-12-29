#!/usr/bin/perl
##---------------------------------------------------------------------------##
##  File:
##      @(#) man.cgi 1.2 97/08/12 12:58:26 @(#)
##  Author:
##      Earl Hood       ehood@medusa.acs.uci.edu
##  Description:
##	man.cgi is a CGI program for viewing Unix manpages.  The
##	program man2html,
##	<URL:http://www.oac.uci.edu/indiv/ehood/man2html.html>,
##	is used to convert the output from man(1) to html.
##
##	If man.cgi is invoked with no input data, it will output a
##	form for the user to select a manpage to view.
##	man.cgi can handle POST and GET methods.
##
##	The code section "Configureable Globals" is designed to
##	allow you to modify man.cgi to work with your particular
##	system configuration.
##---------------------------------------------------------------------------##
##  Copyright (C) 1995-1997,  Earl Hood, ehood@medusa.acs.uci.edu
##
##  This program is free software; you can redistribute it and/or modify
##  it under the terms of the GNU General Public License as published by
##  the Free Software Foundation; either version 2 of the License, or
##  (at your option) any later version.
##  
##  This program is distributed in the hope that it will be useful,
##  but WITHOUT ANY WARRANTY; without even the implied warranty of
##  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
##  GNU General Public License for more details.
##  
##  You should have received a copy of the GNU General Public License
##  along with this program; if not, write to the Free Software
##  Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA
##  02111-1307, USA
##---------------------------------------------------------------------------##

########################################################################
##	Configureable Globals
########################################################################
##	Change the value of these variables to reflect your
##	system configuration.
########################################################################

##  "English" name for program

$ProgName	= "Manpage Viewer";

##  cgi-bin directory as accessed by a URL.

$CgiBin  	= "/cgi-bin";

##  man program (should be a full pathname)

$ManPrg		= '/usr/bin/man';

##  man(1) option to specify a section.  "-s" is used for Solaris
##  systems.  If using Tom Christiansen's Perl man program, than
##  define use "".

#$SectOpt	= "";
$SectOpt	= "-s";

##  man2html program (needs to be a full pathname)

$ManConvPrg	= '/usr/bin/man2html';

##  Flag if the -cgiurl option should be used

$DoCgiurl	= 1;

##  System specific arguments to man2html:
##	HP  => ("-leftm", "1", "-topm", "8")
##	Sun => ("-sun")
##  See man2html documentation for more information.

#@ConvArgs	= ("-compress");
#@ConvArgs	= ("-compress", "-leftm", "1", "-topm", "8");
@ConvArgs	= ("-compress", "-sun");

##  Keyword search processing arguments for man2html.  Normally,
##  '-k' is adequate.  However, if on a Solaris system, the
##  '-solaris' option should be specified with '-k'.  See
##  man2html documentation on information about the '-solaris' option.

#@KeyWArgs	= ("-k");			# Normal
@KeyWArgs	= ("-k", "-solaris");		# For Solaris

##  Possible manual sections.  This array is used to determine the
##  the choices available in an option menu.

@Sections	= (
    '1', '1F', '1M',
    '2',
    '3', '3C', '3F', '3G', '3I', '3N', '3S', '3X',
    '4',
    '5',
    '6',
    '7',
    '8',
    '9', '9F',
);

##  Form method.  The value is either 'GET' or 'POST'.  'GET' is
##  recommended since the URL sent by the client also contains
##  the argument information.  This allows a client's "Reload" function
##  to reprocess a currently viewed manpage.

$FormMethod	= 'GET';

##  Argument separator for CGI URL links.  As clients become more
##  SGML conformant, the simple use of '&' conflicts with
##  SGML syntax.  You can set this variable to control what is
##  used as the separator.  Possibilities:
##	&amp;
##	&#38;
##	;
##

$ArgSep		= '&';

##  Man directories.  Add paths to the list you want man(1) to
##  know about

@ManPath	= qw(
    /usr/man
    /usr/openwin/man
    /usr/man
);

##  PATH setting.  Modify as see fit.  Once useful modification
##  is to have groff utils first in path since its nroff would
##  be invoked over the systems nroff when man formats a manpage.

@Path   	= qw(
    /opt/FSFgroff/bin
    /bin
    /usr/bin
);

########################################################################
##	END!!!! Configureable Globals section
########################################################################

########################################################################
##	Globals
########################################################################

($PROG = $0)	=~ s/.*\///;		# Name of program
$VERSION	= '2.0.1';		# Version
%FORM		= ();			# Hash to hold form contents
$Error		= '';			# Error string
$ENV{'MANPATH'}	= join(":", @ManPath);
$ENV{'PATH'}	= join(":", @Path);

########################################################################
##	Main block
{
    #	Set unbuffered I/O.  Prevents buffering problems with
    #	"system()" calls.
    select((select(STDOUT), $| = 1)[0]);

    #	Print content-type header
    printouttype("text/html");

    #	Print man form if called w/no arguments
    printform() if noarg();

    #	If reached here, there is input to process
    error("CGI input error") unless parseinput();
    printform() unless $FORM{'topic'};
    doit();
    exit 0;
}
########################################################################
##	Subroutines
########################################################################

#-----------------------------------------------------------------------
#	printform outputs the man selection form to the client.
#
sub printform {
    printhead($ProgName);
    print STDOUT <<EndOfForm;
<p>The following form allows you to view a manpage on this system.
Please fill out the following fields and select <strong>Submit</strong>
to view a manpage.
</p>
<form method="$FormMethod" action="$CgiBin/man.cgi">
<table border=0>
<tr>
<td align=right>Section:</td>
<td><select name=section>
<option value="all">All Sections</option>
<option value="keyword">Keyword Search</option>
EndOfForm

    #	Print out possible section choices
    local($section);
    foreach $section (@Sections) {
	print STDOUT qq|<option value="$section">Section $section</option>\n|;
    }

print STDOUT <<EndOfForm;
</select>
</td>
</tr>
<tr>
<td align=right>Topic:</td>
<td><input type="TEXT" name="topic"></td>
</tr>
<tr>
<td></td><td><input type="SUBMIT" value="Submit"><td>
</tr>
</table>
</form>
EndOfForm
    printend();
    exit 0;
}

#-----------------------------------------------------------------------
#	doit does the conversion
#
sub doit {
    my($section, $topic, $manexec, $manout, $tmp);
    $manout = '';

    # Get section and topic from input
    #---------------------------------
    $section = $FORM{'section'};
    $topic = $FORM{'topic'};
    error("Questionable characters in topic")  if isquestionable($topic);

    # Determine command arguments for man and man2html
    #-------------------------------------------------
    @ARGV = ();
    @manargs = ();
    $manexec = $ManPrg;
    if ($section =~ /keyword/) {
	$manexec .= " -k $topic";
	push(@manargs, "-k", $topic);
	push(@ARGV, @KeyWArgs,
		    "-title", qq{Keyword search: "$topic"});
    } else {
	error("No topic entered")  unless $topic;
	if ($section !~ /all/) {
	    push(@manargs, $SectOpt)  if $SectOpt;
	    push(@manargs, $section);
	}
	push(@manargs, $topic);
	$manexec .= " $section"  if $section !~ /all/;
	$manexec .= " $topic";
	$tmp  = $topic;
	$tmp .= "($section)"  if $section !~ /all/;
	push(@ARGV, @ConvArgs,
		    "-title", $tmp);
    }

    # Check if doing man xref detection
    #----------------------------------
    if ($DoCgiurl) {
	push(@ARGV, "-cgiurl",
		    join('', $CgiBin, '/man.cgi?',
			     'section=${section}${subsection}',
			     $ArgSep,
			     'topic=${title}'));
    }

    # Convert output from man to html
    #--------------------------------
    close(STDERR);
    open(MANPRG, "-|") or exec($ManPrg, @manargs);
    $Man2Html::InFH = \*MANPRG;		# set input filehandle
    require $ManConvPrg or
	error("Problem executing man->HTML converter");
}

########################################################################
##	Generic subroutines for CGI use
########################################################################

#-----------------------------------------------------------------------
#	noarg returns true if no arguments were passed to script.
#
sub noarg {
    $ENV{"REQUEST_METHOD"} eq "GET" && $ENV{"QUERY_STRING"} =~ /^\s*$/;
}

#-----------------------------------------------------------------------
#	parseinput converts the input data into the %FORM array
#
sub parseinput {
    my($method) = ($ENV{"REQUEST_METHOD"});
    my($data);
    if ($method eq "GET") {
	$data = $ENV{"QUERY_STRING"};
    } elsif ($method eq "POST") {
	read(STDIN, $data, $ENV{"CONTENT_LENGTH"});
    } else {
	$Error = "Unrecgonized request method : $method";
	return 0;
    }
    my(@pairs, $name, $value);
    if ($data ne '') {
	@pairs = split(/&/, $data);
	foreach (@pairs) {
	    ($name, $value) = split(/=/);
	    $name = expandstr($name);
	    $value = expandstr($value);
	    $FORM{$name} = $value;
	}
    }
    1;
}

#-----------------------------------------------------------------------
#	printouttype prints out specified content-type header back
#	to client
#
sub printouttype {
    my($type) = shift;
    print STDOUT "Content-type: $type\r\n\r\n";
}

#-----------------------------------------------------------------------
#	printhead outputs html prematter
#
sub printhead {
    my($title, $h1) = @_;
    $h1 = $title  unless $h1;

    print STDOUT <<ENDOFHEAD;
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML//EN">
<HTML>
<HEAD>
<TITLE>$title</TITLE>
</HEAD>
<BODY>
<H1>$h1</H1>
ENDOFHEAD
}

#-----------------------------------------------------------------------
#	printend outputs html postmatter
#
sub printend {
    print STDOUT <<ENDOFEND;
<HR>
<ADDRESS>
Manpage viewer available with the
<a href="http://www.oac.uci.edu/indiv/ehood/man2html.html">man2html</a>
package.
</ADDRESS>
</BODY>
</HTML>
ENDOFEND
}

#-----------------------------------------------------------------------
#	error prints an error out to the client.
#
sub error {
    my($str) = htmlize(shift);
    printhead("$ProgName Error");
    $str .= ":"  if $Error && $str;
    $str .= " $Error";
    print STDOUT "<p>$str</p>";
    printend();
    exit 0;
}

#-----------------------------------------------------------------------
#	htmlize translates special characters to enitity refs.
#
sub htmlize {
    my($str) = shift;
    $str =~ s/&/\&amp;/g;
    $str =~ s/</\&lt;/g;
    $str =~ s/>/\&gt;/g;
    $str;
}

#-----------------------------------------------------------------------
#	expandstr translates hex codes to characters
#
sub expandstr {
    my($str) = shift;
    $str =~ tr/+/ /;
    $str =~ s/%([a-fA-F0-9][a-fA-F0-9])/pack("C", hex($1))/ge;
    $str;
}

#-----------------------------------------------------------------------
#	isquestionable determines if $str contains questionable
#	characters if $str is used in a subshell invocation.
#
sub isquestionable {
    my($str) = shift;
    $str !~ /^[a-zA-Z0-9_\-+ \t\/@%\.]+$/;
}

########################################################################
