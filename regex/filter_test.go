package regex

import "testing"

func TestFilter(t *testing.T) {
	testObj := []struct{expr string
						files []string
						}{
		{
			expr:".*exe$",
			files:[]string{"test.exe"},
		},
		{
			expr:`[a-zA-Z0-9]\.[conf|config]`,
			files:[]string{"abc.conf","a12c@.config","Az23.conf"},
		},
	}

	for _,obj:=range testObj{
		for _,file:=range obj.files{
			if Filter(obj.expr,file)!=true{
				t.Errorf("%s is match %s\n",file,obj.expr)
			}
		}

	}
}
