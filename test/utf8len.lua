local text="あいうえおかきくけこさしすせそ"

for i=1,string.len(text) do
    for j=i,string.len(text) do
        local len,err =utf8.len(text,i,j)
        if len then
            print(string.format("[%d:%d]=%d",i,j,len))
        else
            print(string.format("[%d:%d]=ERROR %s",i,j,err))
        end
    end
end
