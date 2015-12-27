package me.wener.bbvm.vm;

import com.google.common.base.Joiner;
import com.google.common.base.Splitter;
import com.google.common.collect.Lists;

import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.Assert.assertEquals;

/**
 * 从 asm 中抽取输入输出
 *
 * @author wener
 * @since 15/12/15
 */
public class TestSpec {
    public static final Pattern REG_MARKER = Pattern.compile(";\\s*(\\|?[><])(.+)");
    public static final Splitter NL_SPLITTER = Splitter.on('\n');
    public static final Joiner NL_JOINER = Joiner.on('\n');
    private final StringBuilder input = new StringBuilder();
    private final StringBuilder output = new StringBuilder();

    public TestSpec clear() {
        input.setLength(0);
        output.setLength(0);
        return this;
    }

    public TestSpec accept(String c) {
        StringBuilder in = new StringBuilder();
        StringBuilder out = new StringBuilder();
        Matcher matcher = REG_MARKER.matcher(c);
        while (matcher.find()) {
            String type = matcher.group(1);
            String content = matcher.group(2).trim();
            if (content.length() == 0) {
                continue;
            }
            switch (type) {
                case "<":
                    in.append(content).append('\n');
                    break;
                case "|<":
                    in.append(content);
                    break;
                case ">":
                    out.append(content).append('\n');
                    break;
                case "|>":
                    out.append(content);
                    break;
                default:
                    throw new AssertionError();
            }
        }
        input(in);
        output(out);
        return this;
    }

    public TestSpec input(CharSequence c) {
        input.append(c);
        return this;
    }

    public TestSpec output(CharSequence c) {
        output.append(c);
        return this;
    }

    public CharSequence output() {
        return output;
    }

    public CharSequence input() {
        return input;
    }

    public void assertMatch(String out) {
        List<Integer> skipped = Lists.newArrayList();
        List<String> expected = NL_SPLITTER.splitToList(output);
        for (int i = 0; i < expected.size(); i++) {
            if (expected.get(i).equals("skip")) {
                skipped.add(i);
            }
        }
        List<String> actually = Lists.newArrayList(NL_SPLITTER.splitToList(out));

        try {
            assertEquals(expected.size(), actually.size());
        } catch (Throwable e) {
            System.out.printf("Expected\n%s\nActually\n%s\n", output, out);
            throw e;
        }

        skipped.forEach(i -> actually.set(i, "skip"));
        assertThat(NL_JOINER.join(expected)).isEqualTo(NL_JOINER.join(actually));
    }
}
