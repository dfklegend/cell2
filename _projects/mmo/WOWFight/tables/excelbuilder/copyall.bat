
set TAR=..\..\..\client\UData\csv\
xcopy /Y ..\export\*.csv %TAR%
del %TAR%\s_*.csv

set TAR=..\..\..\release\client\UData\csv\
xcopy /Y ..\export\*.csv %TAR%
del %TAR%\s_*.csv

set TAR=..\..\..\server\data\csv\
xcopy /Y ..\export\*.csv %TAR%

set TAR=..\..\..\release\server\data\csv\
xcopy /Y ..\export\*.csv %TAR%

rem pause


